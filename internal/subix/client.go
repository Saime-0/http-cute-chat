package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
	"time"
)

type Users map[ID]*User

type User struct {
	ID        int
	membering Members
	clients   Clients
}

type Clients map[Key]*Client
type ClientsWithEvents map[Key]*ClientWithEvents

type Client struct {
	UserID int
	//ExpectedEvents   map[model.EventType]bool
	Ch               chan *model.SubscriptionBody
	task             *scheduler.Task
	sessionExpiresAt int64
	sessionKey       Key
	marked           bool
}

type EventCollection map[model.EventType]bool
type ClientWithEvents struct {
	Client *Client
	Events EventCollection
}

func (s *Subix) CreateUserIfNotExists(userID int) *User {
	user, ok := s.users[userID]
	if !ok {
		user = &User{
			ID:        userID,
			membering: Members{},
			clients:   Clients{},
		}
		s.users[userID] = user
	}
	return user
}

func (s *Subix) deleteUser(userID int) {
	user, ok := s.users[userID]
	if ok { // если вдруг не удается найти то просто скипаем
		delete(s.users, userID)               // удаление из глобальной мапы
		for _, client := range user.clients { // определяем тех клиентов которых надо удалить отовсюду
			delete(s.clients, client.sessionKey) // удлаение из глобальной мапы
		}
		user.clients = nil // на всякий случай заnullяем мапу

		for _, member := range user.membering { // а здесь определяем мемберов, к которые относятся к пользователю
			s.DeleteMember(member.ID) // удаляем по отдельности через функцию
		}
		user.membering = nil // на всякий случай заnullяем мапу
		// теперь на этого пользователя не должно остаться ссылок как и на его клиентов
	}

}

func (s *Subix) deleteClient(sessionKey Key) error {
	client, ok := s.clients[sessionKey]
	if ok {
		delete(s.clients, client.sessionKey)
		err := s.sched.DropTask(&client.task)
		if err != nil {
			return cerrors.Wrap(err, "не удалось удалить клиента")
		}
		close(client.Ch)

		user, ok := s.users[client.UserID]
		if ok {
			delete(user.clients, client.sessionKey)
			if len(user.clients) == 0 {
				s.deleteUser(user.ID)
			}
		}

		for _, member := range user.membering {
			delete(member.clientsWithEvents, client.sessionKey)
			if len(member.clientsWithEvents) == 0 {
				s.DeleteMember(member.ID)
			}
		}
	}
	//println("удален клиент", client.sessionKey)
	return nil
}

func (s *Subix) scheduleMarkClient(client *Client, expAt int64) (err error) {
	client.task, err = s.sched.AddTask(
		func() {
			eventBody := &model.TokenExpired{
				Message: "используйте mutation.RefreshTokens для того чтобы возобновить получение данных, иначе соединение закроется",
			}
			s.writeToClient(
				client,
				&model.SubscriptionBody{
					Event: getEventTypeByEventResult(eventBody),
					Body:  eventBody,
				},
			)
			client.marked = true // теперь будем знать что этому клиенту не надо отправлять события
			//println("токен клиента истек, помечаю клиента", client)
			err := s.scheduleExpiredClient(client)
			if err != nil {
				panic(err)
			}

		},
		expAt,
	)

	return err
}

func (s *Subix) scheduleExpiredClient(client *Client) (err error) {
	client.task, err = s.sched.AddTask(
		func() {
			//fmt.Printf("клиент %s не обновил токен, удаляю", client.sessionKey)
			s.deleteClient(client.sessionKey)
		},
		time.Now().Unix()+rules.LifetimeOfMarkedClient,
	)

	return err
}

func (s *Subix) ExtendClientSession(sessionKey Key, expAt int64) (err error) {
	client, ok := s.clients[sessionKey]
	if !ok {
		return cerrors.New("не удалось продлить сессию, клиент не найден")
	}
	err = s.sched.DropTask(&client.task)
	if err != nil {
		return err
	}
	err = s.scheduleMarkClient(client, expAt)
	if err != nil {
		return err
	}
	client.marked = false
	//println("сессия продлена клиента", client)
	return nil
}

func (s Subix) ClientCollection(sessionKey Key) (collection []*model.ListenedChat) {
	client, _ := s.clients[sessionKey]                        // предполагается что сессия с таким ключем существует
	for _, member := range s.users[client.UserID].membering { // по мемберсам пользователя
		listenedChat := &model.ListenedChat{ID: member.ChatID}
		for event, _ := range member.clientsWithEvents[sessionKey].Events {
			listenedChat.Events = append(listenedChat.Events, event)
		}
		collection = append(collection, listenedChat)
	}
	return collection
}

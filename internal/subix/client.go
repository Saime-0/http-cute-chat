package subix

import (
	"errors"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
	"time"
)

type Users map[ID]*User

type User struct {
	ID         int
	rootClient *Client
}

type Clients map[Key]*Client

type Client struct {
	Ch               chan *model.SubscriptionBody
	UserID           int
	next             **Client
	task             *scheduler.Task
	sessionExpiresAt int64
	webKey           Key
	marked           bool
}

func (s *Subix) CreateUserIfNotExists(userID int) *User {
	user, ok := s.users[userID]
	if !ok {
		user = &User{
			ID:         userID,
			rootClient: &Client{},
		}
		s.users[userID] = user
		println("Создан user id", userID) // debug
	}
	return user
}

func (s *Subix) deleteUser(userID int) {
	delete(s.users, userID)
	println("удален пользователь с  id =", userID) // debug
	s.deleteMemberByUserID(userID)
}

func (s *Subix) deleteClient(client **Client) {
	fmt.Printf("deleteClient: %#v\n", *client) // debug
	delete(s.clients, (*client).webKey)
	err := s.sched.DropTask(&(*client).task)
	if err != nil {
		println("deleteClient:", err.Error())
		return
	}
	close((*client).Ch)               // закрываем канал, потому как функция может вызываться а не триггериться закрытием соединения
	println("удален клиент", *client) // debug
	if (*client).next == nil {
		s.deleteUser((*client).UserID)
		return
	}
	*client = *(*client).next
}

func (s *Subix) scheduleMarkClient(client *Client, expAt int64) (err error) {
	(*client).task, err = s.sched.AddTask(
		func() {
			eventBody := &model.TokenExpired{
				Message: "используйте mutation.RefreshTokens для того чтобы возобновить получение данных, иначе соединение закроется",
			}
			s.writeToClient(
				&client,
				&model.SubscriptionBody{
					Event: getEventType(eventBody),
					Body:  eventBody,
				},
			)
			client.marked = true                                    // теперь будем знать что этому клиенту не надо отправлять события
			println("токен клиента истек, помечаю клиента", client) // debug
			err := s.scheduleExpiredClient(client)
			if err != nil {
				panic(err)
			}

		},
		expAt,
	)
	if err != nil {
		println("scheduleMarkClient", err) // debug
	}
	return err
}

func (s *Subix) scheduleExpiredClient(client *Client) (err error) {
	client.task, err = s.sched.AddTask(
		func() {
			println("клиент не обновил токен, удаляю", client) // debug
			s.deleteClient(&client)
		},
		time.Now().Unix()+rules.LifetimeOfMarkedClient,
	)
	if err != nil {
		println("scheduleMarkClient", err) // debug
	}
	return err
}

func (s *Subix) ExtendClientSession(webKey Key, expAt int64) (err error) {
	client, ok := s.clients[webKey]
	if !ok {
		return errors.New("не удалось продлить сессию, клиент не найден")
	}
	err = s.sched.DropTask(&client.task)
	if err != nil {
		println("extendlientSession", err) // debug
		return err
	}
	err = s.scheduleMarkClient(client, expAt)
	if err != nil {
		println("extendlientSession", err) // debug
		return err
	}
	client.marked = false
	println("сессия продлена клиента", client) // debug
	return nil
}

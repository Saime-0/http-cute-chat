package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (s *Subix) Sub(userID int, sessionKey Key, expAt int64) (*Client, error) {
	_, ok := s.clients[sessionKey]
	if ok { // если ключ сущесивует то по хорошему клиент должен повторить соединение с другим ключем
		return nil, cerrors.New("sessionKey already in use, it is not possible to create a new connection")
	}

	// websocket connection = сессия = клиент
	client := &Client{
		UserID:           userID,
		Ch:               make(chan *model.SubscriptionBody),
		sessionExpiresAt: expAt,
		sessionKey:       sessionKey,
	}
	s.clients[sessionKey] = client

	// планируем пометку и дальнейшее удаление клиента, если его токен истечет
	err := s.scheduleMarkClient(client, expAt)
	if err != nil {
		delete(s.clients, sessionKey)
		return nil, cerrors.New("не удалось создать сессию")
	}

	// user
	user := s.CreateUserIfNotExists(userID)
	user.clients[sessionKey] = client

	return client, nil
}

var allUsefulEventtypes = model.AllEventType[1:len(model.AllEventType)]

func (s *Subix) ModifyCollection(sessionKey Key, submembers []*models.SubUser, action model.EventSubjectAction, listenEvents []model.EventType) error {
	client, ok := s.clients[sessionKey]
	if !ok { // если ключ сущесивует, то по-хорошему клиент должен повторить соединение с другим ключем
		return cerrors.New("no session with this key was found")
	}
	user := s.users[client.UserID]
	for _, event := range listenEvents {
		if event == model.EventTypeAll {
			listenEvents = allUsefulEventtypes
			break
		}
	}
	if action == model.EventSubjectActionAdd { // если мемберс добавляет пачку событий
		for _, sm := range submembers {
			member := s.CreateMemberIfNotExists(*sm.MemberID, *sm.ChatID, user.ID) // достаем мембера из активных(те на которые подписаны клиенты) нужного мембера
			clientWithEvents, ok := member.clientsWithEvents[sessionKey]           // ищем сессию нужного клиента в мемберсе
			if !ok {                                                               // если клиент еще не прорслушивает этого участника, то заставляем сушать
				clientWithEvents = &ClientWithEvents{
					Client: client,
					Events: make(EventCollection),
				}
				member.clientsWithEvents[sessionKey] = clientWithEvents
			}
			for _, event := range listenEvents {
				clientWithEvents.Events[event] = true // добавил тип ивента который теперь будет отправляться клиенту(прослушиваться им)
			}

			// add member to user memberings
			user.membering[*sm.MemberID] = member // даже если у пользователя уже есть мембер с таким id то все равно добавляем(не даст никакого эффекта)
		}
	}

	if action == model.EventSubjectActionDelete {
		for _, sm := range submembers {

			member, ok := s.members[*sm.MemberID]
			if !ok {
				continue // пропускаем если клиент не слушает этого участника
			} else {
				clientWithEvents, ok := member.clientsWithEvents[sessionKey]
				if !ok {
					break
				}
				for _, event := range listenEvents {
					delete(clientWithEvents.Events, event)
				}

				if len(clientWithEvents.Events) == 0 { // если удалятся все события которые прослушывал клиент, то ..
					delete(member.clientsWithEvents, sessionKey) // .. удаляем клиента из мемберса
					if len(member.clientsWithEvents) == 0 {      // а если количесво слушающих клиентов = 0 то..
						s.DeleteMember(*sm.MemberID) // удаляем мембера
					}
				}
			}

		}
	}
	return nil
}

func (s *Subix) Unsub(sessionKey Key) error {
	err := s.deleteClient(sessionKey)
	if err != nil {
		return cerrors.Wrap(err, "не удалось отписаться")
	}
	return nil
}

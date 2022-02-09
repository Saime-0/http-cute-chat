package subix

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (s *Subix) Sub(userID int, sessionKey Key, expAt int64, submembers []*models.SubUser) (*Client, error) {
	_, ok := s.clients[sessionKey]
	if ok { // если ключ сущесивует то по хорошему клиент должен повторить соединение с другим ключем
		return nil, errors.New("sessionKey already in use, it is not possible to create a new connection")
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
		return nil, errors.New("не удалось создать сессию")
	}

	// user
	user := s.CreateUserIfNotExists(userID)
	user.clients[sessionKey] = client

	for _, sm := range submembers {
		member := s.CreateMemberIfNotExists(*sm.MemberID, *sm.ChatID, user.ID)
		member.clients[sessionKey] = client

		chat := s.CreateChatIfNotExists(*sm.ChatID)
		chat.members[*sm.MemberID] = member

		// add member to user memberings
		user.membering[*sm.MemberID] = member
	}
	return client, nil
}

func (s *Subix) Unsub(sessionKey Key) error {
	err := s.deleteClient(sessionKey)
	if err != nil {
		return errors.Wrap(err, "не удалось отписаться")
	}
	return nil
}

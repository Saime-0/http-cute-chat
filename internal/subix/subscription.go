package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (s *Subix) Sub(userID int, submembers []*models.SubUser) *Client {
	// websocket connection
	client := &Client{
		Ch:     make(chan *model.SubscriptionBody),
		UserID: userID,
		next:   nil,
	}
	println("Создан client ", &client) // debug

	// user
	user := s.CreateUserIfNotExists(userID)
	user.rootClient.next, client.next = &client, user.rootClient.next // добавили новую "сессию" пользователя

	for _, submember := range submembers {
		// chat
		chat := s.CreateChatIfNotExists(*submember.ChatID)

		// member
		s.CreateMemberIfNotExists(*submember.MemberID, user, chat)
	}
	return client
}

func (s *Subix) Unsub(client **Client) {
	if client == nil || *client == nil {
		println("не удалось отписать клиента, тк его не существует") // debug
		return
	}
	s.deleteClient(client)
}

package subix

import (
	"errors"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (s *Subix) Sub(userID int, webSocketKey Key, expAt int64, submembers []*models.SubUser) (*Client, error) {
	if len(webSocketKey) == 0 {
		return nil, errors.New("not valid \"WebSocket-Session-Key\"")
	}
	fmt.Printf("%#v", s.clients) // debug
	_, ok := s.clients[webSocketKey]
	if ok { // если ключ сущесивует то по хорошему клиент должен повторить соединение с другим ключем
		return nil, errors.New("webSocketKey already in use, it is not possible to create a new connection")
	}

	// websocket connection = сессия = клиент
	client := &Client{
		Ch:               make(chan *model.SubscriptionBody),
		UserID:           userID,
		webKey:           webSocketKey,
		sessionExpiresAt: expAt,
	}
	s.clients[webSocketKey] = client
	fmt.Printf("%#v", client) // debug

	// планируем пометку и дальнейшее удаление клиента, если его токен истечет
	err := s.scheduleMarkClient(client, expAt)
	if err != nil {
		return nil, errors.New("не удалось создать сессию")
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
	return client, nil
}

func (s *Subix) Unsub(client **Client) {
	fmt.Printf("Клиент хочет удалиться %#v\n", *client) // debug
	if client == nil || *client == nil {
		println("не удалось отписать клиента, тк его не существует") // debug
		return
	}
	s.deleteClient(client)
}

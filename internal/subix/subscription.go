package subix

import (
	"errors"
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
	println("Создан клиент", sessionKey) // debug

	// user
	user := s.CreateUserIfNotExists(userID)
	user.clients[sessionKey] = client

	for _, sm := range submembers {
		// member
		//member := func(memberID, chatID int) *Member {
		//	member, ok := s.members[memberID]
		//	if !ok {
		//		member = &Member{
		//			ID:     memberID,
		//			ChatID: chatID,
		//			UserID: userID,
		//		}
		//		s.members[memberID] = member
		//	}
		//	return member
		//}(*sm.MemberID, *sm.ChatID)
		member := s.CreateMemberIfNotExists(*sm.MemberID, *sm.ChatID, user.ID)
		member.clients[sessionKey] = client

		// chat
		//chat := func(chatID int) *Chat {
		//	chat, ok := s.chats[chatID]
		//	if !ok {
		//		chat = &Chat{
		//			ID:      0,
		//			members: Members{},
		//		}
		//		s.chats[chatID] = chat
		//		println("Создан chat id", chat.ID) // debug
		//	}
		//	return chat
		//}(*sm.ChatID)
		chat := s.CreateChatIfNotExists(*sm.ChatID)
		chat.members[*sm.MemberID] = member

		// add member to user memberings
		user.membering[*sm.MemberID] = member
	}
	return client, nil
}

func (s *Subix) Unsub(sessionKey Key) {
	println("Клиент хочет удалиться", sessionKey) // debug
	s.deleteClient(sessionKey)
}

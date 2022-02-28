package subix

type Chats map[ID]*Chat

type Chat struct {
	ID      int
	members Members
}

type Members map[ID]*Member

type Member struct {
	ID                int
	ChatID            int
	UserID            int
	clientsWithEvents ClientsWithEvents
}

func (s *Subix) CreateChatIfNotExists(chatID int) *Chat {
	chat, ok := s.chats[chatID]
	if !ok {
		chat = &Chat{
			ID:      chatID,
			members: Members{},
		}
		s.chats[chatID] = chat
	}
	return chat
}

func (s *Subix) CreateMemberIfNotExists(memberID, chatID, userID int) *Member {
	member, ok := s.members[memberID]
	if !ok {
		chat := s.CreateChatIfNotExists(chatID)
		chat.members[memberID] = member
		member = &Member{
			ID:                memberID,
			ChatID:            chatID,
			UserID:            userID,
			clientsWithEvents: make(ClientsWithEvents),
		}

		s.members[memberID] = member
	}
	return member
}

func (s *Subix) deleteChat(chatID int) {
	chat, ok := s.chats[chatID]
	if ok {

		for _, member := range chat.members {
			s.DeleteMember(member.ID)
		}
		delete(s.chats, chatID)
	}
}

func (s *Subix) DeleteMember(memberID int) {
	member, ok := s.members[memberID]
	if ok { // если вдруг не удается найти то просто скипаем
		delete(s.members, memberID)    // удлаение из глобальной мапы
		member.clientsWithEvents = nil // на всякий случай заnullяем мапу

		user, ok := s.users[member.UserID]
		if ok {
			delete(user.membering, member.ID)
			//if len(user.membering) == 0 { это не должно здесь быть
			//	s.deleteUser(user.ID)
			//}
		}

		chat, ok := s.chats[member.ChatID]
		if ok {
			delete(chat.members, member.ID)
			if len(chat.members) == 0 {
				s.deleteChat(chat.ID)
			}
		}
	}
}

package subix

import "fmt"

type Chats map[ID]*Chat

type Chat struct {
	ID         int
	rootMember *Member
}

type Members map[ID]**Member

type Member struct {
	ID     int
	ChatID int
	User   *User
	next   **Member
}

func (s *Subix) CreateChatIfNotExists(chatID int) *Chat {
	chat, ok := s.chats[chatID]
	if !ok {
		chat = &Chat{
			ID:         chatID,
			rootMember: &Member{},
		}
		s.chats[chatID] = chat
		println("Создан chat id", chat.ID) // debug
	}
	return chat
}

func (s *Subix) CreateMemberIfNotExists(memberID int, user *User, chat *Chat) **Member {
	member, ok := s.members[memberID]
	if !ok {
		_member := &Member{
			ID:     memberID,
			ChatID: chat.ID,
			User:   user,
			next:   nil,
		}
		member = &_member
		s.members[memberID] = member
		// если такого мембера еще небыло создано, те до этого созданной сессии пользователя, то мы создаем, а если есть то пофиг, тк пользователь бует уже привязан
		chat.rootMember.next, (*member).next = member, chat.rootMember.next
		println("Создан member id", (*member).ID) // debug
	}
	return member
}

func (s *Subix) deleteChat(chatID int) {
	delete(s.chats, chatID)
	println("удален чат с id =", chatID) // debug
}

func (s *Subix) deleteMemberByUserID(userID int) {
	for _, member := range s.members {
		if (*member).User.ID == userID {
			delete(s.members, (*member).ID)
			fmt.Printf("удален участник чата с id = %d (uid:%d)\n", (*member).ID, userID) // debug
			if (*member).next == nil {
				s.deleteChat((*member).ChatID)
				continue
			}
			*member = *(*member).next
		}
	}
}

func (s *Subix) deleteMember(memberID int) {
	member, ok := s.members[memberID]
	if !ok {
		println("не удалось удалить участника, участник с id =", memberID, "не найден")
		return
	}
	*member = *(*member).next
	delete(s.members, memberID)
	println("удален участник чата с id =", memberID) // debug
}

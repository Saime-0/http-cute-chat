package piping

import "github.com/saime-0/http-cute-chat/internal/its"

func (p *Pipeline) ChatExists(chatId int) (fail bool) {
	return !p.Repos.Chats.ChatExistsByID(chatId)
}

func (p *Pipeline) UserExists(userId int) (fail bool) {
	return !p.Repos.Users.UserExistsByID(userId)
}

func (p *Pipeline) UserIs(chatId, userId int, somes []its.Someone) (fail bool) {
	for _, some := range somes {
		switch some {
		case its.Owner:
			return !p.Repos.Chats.UserIsChatOwner(userId, chatId)

		case its.Admin:
			//p.Repos.Chats.
			return

		case its.Moder:
			//p.Repos.Chats.
			return

		case its.Member:
			return !p.Repos.Chats.UserIsChatMember(userId, chatId)

		}
	}

	return
}

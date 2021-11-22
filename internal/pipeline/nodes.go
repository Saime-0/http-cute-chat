package pipeline

import "github.com/saime-0/http-cute-chat/internal/its"

func (p Pipeline) ChatExists() (fail bool) {
	return
}

func (p Pipeline) UserExists() (fail bool) {
	return
}

func (p Pipeline) UserIs(somes ...its.Someone) (fail bool) {
	// пользователь должен будет сам определяться, через какой нибудь метод
	// который возвращает либо target либо просто user, а это поля переданы поле контекста конструктора пайплайна
	// либо проинициализированы сущности
	user_id := p.Ctx.Value("user_id").(int)
	chat_id := p.Ctx.Value("chat_id").(int)
	for _, some := range somes {
		switch some {
		case its.Owner:
			return !p.Repos.Chats.UserIsChatOwner(user_id, chat_id)

		case its.Admin:
			//p.Repos.Chats.
			return

		case its.Moder:
			//p.Repos.Chats.
			return

		case its.Member:
			return !p.Repos.Chats.UserIsChatMember(user_id, chat_id)

		}
	}

	return
}

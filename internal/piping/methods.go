package piping

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/api/validator"
	"github.com/saime-0/http-cute-chat/internal/its"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (p *Pipeline) ChatExists(chatId int) (fail bool) {
	if !p.repos.Units.UnitExistsByID(chatId, rules.Chat) {
		p.Err = resp.Error(resp.ErrBadRequest, "такого чата не существует")
		return true
	}
	return
}
func (p *Pipeline) ChatExistsByDomain(chatDomain string) (fail bool) {
	if !p.repos.Units.UnitExistsByDomain(chatDomain, rules.Chat) {
		p.Err = resp.Error(resp.ErrBadRequest, "такого чата не существует")
		return true
	}
	return
}

func (p *Pipeline) UserExists(userId int) (fail bool) {
	if !p.repos.Units.UnitExistsByID(userId, rules.User) {
		p.Err = resp.Error(resp.ErrBadRequest, "пользователь не найден")
		return true
	}
	return
}

func (p *Pipeline) UserIs(chatId, userId int, somes []its.Someone) (fail bool) {
	for _, some := range somes {
		switch some {
		case its.Owner:
			if !p.repos.Chats.UserIsChatOwner(userId, chatId) {
				p.Err = resp.Error(resp.ErrBadRequest, "пользователь не является владельцем чата")
			}
			return

		case its.Admin:
			if !p.repos.Chats.UserIs(userId, chatId, rules.Admin) {
				p.Err = resp.Error(resp.ErrBadRequest, "пользователь не является администратором")
			}
			return

		case its.Moder:
			if !p.repos.Chats.UserIs(userId, chatId, rules.Moder) {
				p.Err = resp.Error(resp.ErrBadRequest, "пользователь не является модератором")
			}
			return

		case its.Member:
			if !p.repos.Chats.UserIsChatMember(userId, chatId) {
				p.Err = resp.Error(resp.ErrBadRequest, "пользователь не является участником чата")
			}
			return

		}
	}

	return
}

// GetIDByDomain
// put ID ptr to 2 arg
func (p *Pipeline) GetIDByDomain(domain string, id *int) (fail bool) {
	_id, err := p.repos.Units.GetIDByDomain(domain)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*id = _id

	return
}

// ValidParams
//  with side effect
func (p *Pipeline) ValidParams(params *model.Params) (fail bool) {
	if params == nil {
		params = &model.Params{
			Limit:  kit.IntPtr(rules.MaxLimit), // ! unsafe
			Offset: kit.IntPtr(0),
		}
	}
	if !validator.ValidateLimit(*params.Limit) {
		p.Err = resp.Error(resp.ErrBadRequest, "невалидное значение лимита")
		return true
	}
	if !validator.ValidateOffset(*params.Offset) {
		p.Err = resp.Error(resp.ErrBadRequest, "невалидное значение смещения")
		return true
	}
	return
}

func (p *Pipeline) ValidNameFragment(fragment string) (fail bool) {
	if validator.ValidateNameFragment(fragment) {
		p.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для фрагмента имени")
		return true
	}
	return
}

func (p *Pipeline) OwnedLimit(userId int) (fail bool) {
	count, err := p.repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count > rules.MaxCountOwnedChats {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит созднных чатов")
		return true
	}
	return
}

func (p *Pipeline) ChatCountLimit(userId int) (fail bool) {
	count, err := p.repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxUserChats {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества чатов в которых пользователь может состоять")
		return true
	}
	return
}

func (p *Pipeline) ValidDomain(domain string) (fail bool) {
	if validator.ValidateDomain(domain) {
		p.Err = resp.Error(resp.ErrBadRequest, "невалидный домен")
		return true
	}
	return
}

func (p *Pipeline) ValidName(name string) (fail bool) {
	if validator.ValidateName(name) {
		p.Err = resp.Error(resp.ErrBadRequest, "невалидное имя")
		return true
	}
	return
}

func (p *Pipeline) DomainIsFree(domain string) (fail bool) {
	if !p.repos.Units.DomainIsFree(domain) {
		p.Err = resp.Error(resp.ErrBadRequest, "домен занят")
		return true
	}
	return
}

func (p *Pipeline) GetUserChar(userId int, chatId int, char *rules.CharType) (fail bool) {
	_char, err := p.repos.Chats.GetUserChar(userId, chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*char = _char
	return
}

func (p *Pipeline) CountInviteLimit(chatId int) (fail bool) {
	count, err := p.repos.Chats.GetCountLinks(chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxInviteLinks {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества инвайтов")
		return true
	}
	return
}

func (p *Pipeline) ValidInviteInput(inp model.CreateInviteInput) (fail bool) {
	if *inp.Exp != 0 && !validator.ValidateLifetime(int64(*inp.Exp)) {
		p.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение времени жизни ссылки")
		return true
	}
	if *inp.Aliens != 0 && !validator.ValidateAliens(*inp.Aliens) {
		p.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение количества ипользований ссылки")
		return true
	}
	return
}

func (p *Pipeline) IsMember(userId, chatId int) (fail bool) {
	if !p.repos.Chats.UserIsChatMember(userId, chatId) {
		p.Err = resp.Error(resp.ErrBadRequest, "пользователь не является участником чата")
		return true
	}
	return
}

func (p *Pipeline) CountRoleLimit(chatId int) (fail bool) {
	count, err := p.repos.Chats.GetCountChatRoles(chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count > rules.MaxRolesInChat {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества ролей в чате")
		return true
	}
	return
}

func (p *Pipeline) CountRoomLimit(chatId int) (fail bool) {
	count, err := p.repos.Chats.GetCountRooms(chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count > rules.MaxCountRooms {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества комнат в чате")
		return true
	}
	return
}

func (p *Pipeline) RoomExists(roomId int) (fail bool) {
	if !p.repos.Rooms.RoomExistsByID(roomId) {
		p.Err = resp.Error(resp.ErrBadRequest, "такой комнаты не существует")
		return true
	}
	return
}

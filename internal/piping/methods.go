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

func (p *Pipeline) GetIDByDomain(domain string, id *int) (fail bool) {
	_id, err := p.repos.Units.GetIDByDomain(domain)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*id = _id

	return
}

// ValidParams with side effect
func (p *Pipeline) ValidParams(params *model.Params) (fail bool) {
	if params == nil {
		params = model.Params{
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

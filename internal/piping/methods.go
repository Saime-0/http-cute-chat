package piping

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/api/validator"
	"github.com/saime-0/http-cute-chat/internal/its"
	"github.com/saime-0/http-cute-chat/internal/models"
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
	if !validator.ValidateNameFragment(fragment) {
		p.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для фрагмента имени")
		return true
	}
	return
}

func (p *Pipeline) ValidID(id int) (fail bool) {
	if !validator.ValidateID(id) {
		p.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для id")
		return true
	}
	return
}

func (p *Pipeline) ValidForm(form *model.UpdateFormInput) (fail bool) {
	_, err := validator.ValidateRoomForm(form)
	if err != nil {
		p.Err = resp.Error(resp.ErrBadRequest, err.Error())
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
	if count >= rules.MaxCountOwnedChats {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит созднных чатов")
		return true
	}
	return
}

func (p *Pipeline) ChatsLimit(userId int) (fail bool) {
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
	if !validator.ValidateDomain(domain) {
		p.Err = resp.Error(resp.ErrBadRequest, "невалидный домен")
		return true
	}
	return
}

func (p *Pipeline) ValidName(name string) (fail bool) {
	if !validator.ValidateName(name) {
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

func (p *Pipeline) InvitesLimit(chatId int) (fail bool) {
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
	if *inp.Duration != 0 && !validator.ValidateLifetime(*inp.Duration) {
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

func (p *Pipeline) IsNotMember(userId, chatId int) (fail bool) {
	if p.repos.Chats.UserIsChatMember(userId, chatId) {
		p.Err = resp.Error(resp.ErrBadRequest, "пользователь является участником чата")
		return true
	}
	return
}

func (p *Pipeline) RolesLimit(chatId int) (fail bool) {
	count, err := p.repos.Chats.GetCountChatRoles(chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxRolesInChat {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества ролей в чате")
		return true
	}
	return
}

func (p *Pipeline) RoomsLimit(chatId int) (fail bool) {
	count, err := p.repos.Chats.GetCountRooms(chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxCountRooms {
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

func (p *Pipeline) IsNotChild(roomId int) (fail bool) {
	if p.repos.Rooms.HasParent(roomId) {
		p.Err = resp.Error(resp.ErrBadRequest, "комната является веткой")
		return true
	}
	return
}

func (p *Pipeline) HasInvite(chatId int, code string) (fail bool) {
	if !p.repos.Chats.HasInvite(chatId, code) {
		p.Err = resp.Error(resp.ErrBadRequest, "такого кода не существует")
		return true
	}
	return
}

func (p *Pipeline) InviteIsRelevant(code string) (fail bool) {
	if !p.repos.Chats.InviteIsRelevant(code) {
		p.Err = resp.Error(resp.ErrBadRequest, "инвайт неактуален")
		return true
	}
	return
}

func (p *Pipeline) RoleExists(chatID, roleID int) (fail bool) {
	if !p.repos.Chats.RoleExistsByID(chatID, roleID) {
		p.Err = resp.Error(resp.ErrBadRequest, "такой роли не существует")
		return true
	}
	return
}
func (p *Pipeline) InviteIsExists(code string) (fail bool) {
	if !p.repos.Chats.InviteExistsByCode(code) {
		p.Err = resp.Error(resp.ErrBadRequest, "инвайта не существует")
		return true
	}
	return
}

func (p *Pipeline) GetChatByInvite(code string, chatId *int) (fail bool) {
	_id, err := p.repos.Chats.ChatIDByInvite(code)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*chatId = _id
	return
}

func (p *Pipeline) MembersLimit(chatId int) (fail bool) {
	count, err := p.repos.Chats.CountMembers(chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxMembersOnChat {
		p.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества участников в чате")
		return true
	}
	return
}

func (p *Pipeline) ChatIsNotPrivate(chatId int) (fail bool) {
	if p.repos.Chats.ChatIsPrivate(chatId) {
		p.Err = resp.Error(resp.ErrBadRequest, "этот чат запривачин")
		return true
	}
	return
}

func (p *Pipeline) UserExistsByInput(input model.LoginInput) (fail bool) {
	if !p.repos.Users.UserExistsByInput(&models.UserInput{
		Email:    input.Email,
		Password: input.Password,
	}) {
		p.Err = resp.Error(resp.ErrBadRequest, "пользователь с такими данными не найден")
		return true
	}
	return
}

func (p *Pipeline) GetUserIDByInput(input model.LoginInput, userId *int) (fail bool) {
	_uid, err := p.repos.Users.GetUserIdByInput(&models.UserInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*userId = _uid
	return
}

func (p *Pipeline) IsNotBanned(userId, chatId int) (fail bool) {
	if p.repos.Chats.UserIsBanned(userId, chatId) {
		p.Err = resp.Error(resp.ErrBadRequest, "вы забанены в этом чате")
		return true
	}
	return
}

func (p *Pipeline) GetChatIDByRoom(roomId int, chatId *int) (fail bool) {
	_chatId, err := p.repos.Rooms.GetChatIDByRoomID(roomId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*chatId = _chatId
	return
}

func (p *Pipeline) MessageAvailable(msgId, roomId int) (fail bool) {
	if !p.repos.Messages.MessageAvailableOnRoom(msgId, roomId) {
		p.Err = resp.Error(resp.ErrBadRequest, "сообщение не найдено")
		return true
	}
	return
}

func (p *Pipeline) IsAllowedTo(action rules.AllowActionType, roomId int, holder *models.AllowHolder) (fail bool) {
	if !p.repos.Rooms.Allowed(action, roomId, holder) {
		p.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав на это действие")
		return true
	}
	return
}

func (p *Pipeline) GetAllowHolder(userId, chatId int, holder *models.AllowHolder) (fail bool) {
	_holder, err := p.repos.Rooms.AllowHolder(userId, chatId)
	if err != nil {
		p.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*holder = *_holder
	return
}

// deprecated
func (p *Pipeline) IsAllowsSet(roomId int) (fail bool) {
	if !p.repos.Rooms.AllowsIsSet(roomId) {
		p.Err = resp.Error(resp.ErrBadRequest, "в комнате не установлены ограничения")
		return true
	}
	return
}

func (p *Pipeline) GetMessageByID(msgId int, message *model.Message) (fail bool) {
	_message, err := p.repos.Messages.Message(msgId)
	if err != nil {
		p.Err = resp.Error(resp.ErrBadRequest, "сообщение не найдено")
		return true
	}
	message = _message
	return
}

func (p *Pipeline) FindMember(memberId int, chatId *int) (fail bool) {
	chatId = p.repos.Chats.ChatIDByMemberID(memberId)
	if chatId == nil {
		p.Err = resp.Error(resp.ErrBadRequest, "участник не найден")
		return true
	}
	return
}

func (p *Pipeline) GetMemberBy(userId, chatId int, memberId *int) (fail bool) {
	by := p.repos.Chats.FindMemberBy(userId, chatId)
	if by == nil {
		p.Err = resp.Error(resp.ErrBadRequest, "не удалось определить участника чата")
		return true
	}
	return
}

func (p *Pipeline) IsNotMuted(memberId int) (fail bool) {
	if p.repos.Chats.MemberIsMuted(memberId) {
		p.Err = resp.Error(resp.ErrBadRequest, "участник чата заглушен")
		return true
	}
	return
}

func (p *Pipeline) GetMemberIDWhich(userId, chatId int, memberId *int) (fail bool) {

}

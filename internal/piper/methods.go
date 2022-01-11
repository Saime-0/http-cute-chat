package piper

import (
	"encoding/json"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/validator"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"strconv"
)

func (n *Node) ChatExists(chatId int) (fail bool) {
	if !n.repos.Units.UnitExistsByID(chatId, rules.Chat) {
		n.Err = resp.Error(resp.ErrBadRequest, "такого чата не существует")
		return true
	}
	return
}
func (n *Node) ChatExistsByDomain(chatDomain string) (fail bool) {
	if !n.repos.Units.UnitExistsByDomain(chatDomain, rules.Chat) {
		n.Err = resp.Error(resp.ErrBadRequest, "такого чата не существует")
		return true
	}
	return
}

func (n *Node) UserExists(userId int) (fail bool) {
	if !n.repos.Units.UnitExistsByID(userId, rules.User) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь не найден")
		return true
	}
	return
}

// ValidParams
//  with side effect
func (n *Node) ValidParams(params **model.Params) (fail bool) {
	if *params == nil {
		*params = &model.Params{
			Limit:  kit.IntPtr(rules.MaxLimit), // ! unsafe
			Offset: kit.IntPtr(0),
		}
		return
	}
	if (*params).Limit != nil {
		if !validator.ValidateLimit(*(*params).Limit) {
			n.Err = resp.Error(resp.ErrBadRequest, "невалидное значение лимита")
			return true
		}
	} else {
		(*params).Limit = kit.IntPtr(rules.MaxLimit)
	}
	if (*params).Offset != nil {
		if !validator.ValidateOffset(*(*params).Offset) {
			n.Err = resp.Error(resp.ErrBadRequest, "невалидное значение смещения")
			return true
		}
	} else {
		(*params).Offset = kit.IntPtr(0)
	}

	return
}

func (n *Node) ValidNameFragment(fragment string) (fail bool) {
	if !validator.ValidateNameFragment(fragment) {
		n.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для фрагмента имени")
		return true
	}
	return
}
func (n *Node) ValidNote(note string) (fail bool) {
	if !validator.ValidateNote(note) {
		n.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для заметки")
		return true
	}
	return
}

func (n *Node) ValidID(id int) (fail bool) {
	if !validator.ValidateID(id) {
		n.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для id")
		return true
	}
	return
}

func (n *Node) ValidFindMessagesInRoom(find *model.FindMessagesInRoom) (fail bool) {
	if find.Count <= 0 ||
		find.Count > rules.MaxMessagesCount ||
		find.Created == model.MessagesCreatedBefore && find.StartMessageID-find.Count < 0 {
		n.Err = resp.Error(resp.ErrBadRequest, "неверное значение количества сообщений")
		return true
	}
	if !validator.ValidateID(find.StartMessageID) {
		n.Err = resp.Error(resp.ErrBadRequest, "неверный ID сообщения")
		return true
	}

	return
}

func (n *Node) ValidPassword(password string) (fail bool) {
	if !validator.ValidatePassword(password) {
		n.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение для пароля")
		return true
	}
	return
}

func (n *Node) ValidForm(form *model.UpdateFormInput) (fail bool) {
	_, err := validator.ValidateRoomForm(form)
	if err != nil {
		n.Err = resp.Error(resp.ErrBadRequest, err.Error())
		return true
	}
	return
}

func (n *Node) OwnedLimit(userId int) (fail bool) {
	count, err := n.repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxCountOwnedChats {
		n.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит созднных чатов")
		return true
	}
	return
}

func (n *Node) ChatsLimit(userId int) (fail bool) {
	count, err := n.repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxUserChats {
		n.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества чатов в которых пользователь может состоять")
		return true
	}
	return
}

func (n *Node) ValidDomain(domain string) (fail bool) {
	if !validator.ValidateDomain(domain) {
		n.Err = resp.Error(resp.ErrBadRequest, "невалидный домен")
		return true
	}
	return
}

func (n *Node) ValidName(name string) (fail bool) {
	if !validator.ValidateName(name) {
		n.Err = resp.Error(resp.ErrBadRequest, "невалидное имя")
		return true
	}
	return
}

func (n *Node) DomainIsFree(domain string) (fail bool) {
	if !n.repos.Units.DomainIsFree(domain) {
		n.Err = resp.Error(resp.ErrBadRequest, "домен занят")
		return true
	}
	return
}

func (n *Node) GetUserChar(userId int, chatId int, char *rules.CharType) (fail bool) {
	_char, err := n.repos.Chats.GetUserChar(userId, chatId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*char = _char
	return
}

func (n *Node) InvitesLimit(chatId int) (fail bool) {
	count, err := n.repos.Chats.GetCountLinks(chatId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxInviteLinks {
		n.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества инвайтов")
		return true
	}
	return
}

func (n *Node) ValidInviteInput(inp model.CreateInviteInput) (fail bool) {
	if inp.Duration != nil && !validator.ValidateLifetime(*inp.Duration) {
		n.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение времени жизни ссылки")
		return true
	}
	if inp.Aliens != nil && !validator.ValidateAliens(*inp.Aliens) {
		n.Err = resp.Error(resp.ErrBadRequest, "недопустимое значение количества ипользований ссылки")
		return true
	}
	return
}

// IsMember does not need if the Can method is used..
func (n *Node) IsMember(userId, chatId int) (fail bool) {
	if !n.repos.Chats.UserIsChatMember(userId, chatId) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь не является участником чата")
		return true
	}

	return
}

func (n *Node) IsNotMember(userId, chatId int) (fail bool) {
	if n.repos.Chats.UserIsChatMember(userId, chatId) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь является участником чата")
		return true
	}
	return
}

func (n *Node) RolesLimit(chatId int) (fail bool) {
	count, err := n.repos.Chats.GetCountChatRoles(chatId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxRolesInChat {
		n.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества ролей в чате")
		return true
	}
	return
}

func (n *Node) RoomsLimit(chatId int) (fail bool) {
	count, err := n.repos.Chats.GetCountRooms(chatId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxCountRooms {
		n.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества комнат в чате")
		return true
	}
	return
}

func (n *Node) RoomExists(roomId int) (fail bool) {
	if !n.repos.Rooms.RoomExistsByID(roomId) {
		n.Err = resp.Error(resp.ErrBadRequest, "такой комнаты не существует")
		return true
	}
	return
}

func (n *Node) IsNotChild(roomId int) (fail bool) {
	if n.repos.Rooms.HasParent(roomId) {
		n.Err = resp.Error(resp.ErrBadRequest, "комната является веткой")
		return true
	}
	return
}

func (n *Node) HasInvite(chatId int, code string) (fail bool) {
	if !n.repos.Chats.HasInvite(chatId, code) {
		n.Err = resp.Error(resp.ErrBadRequest, "такого кода не существует")
		return true
	}
	return
}

func (n *Node) InviteIsRelevant(code string) (fail bool) {
	if !n.repos.Chats.InviteIsRelevant(code) {
		n.Err = resp.Error(resp.ErrBadRequest, "инвайт неактуален")
		return true
	}
	return
}

func (n *Node) RoleExists(chatID, roleID int) (fail bool) {
	if !n.repos.Chats.RoleExistsByID(chatID, roleID) {
		n.Err = resp.Error(resp.ErrBadRequest, "такой роли не существует")
		return true
	}
	return
}
func (n *Node) InviteIsExists(code string) (fail bool) {
	if !n.repos.Chats.InviteExistsByCode(code) {
		n.Err = resp.Error(resp.ErrBadRequest, "инвайта не существует")
		return true
	}
	return
}
func (n *Node) GetChatIDByRole(roleId int, chatId *int) (fail bool) {
	_chatId, err := n.repos.Chats.ChatIDByRoleID(roleId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "не удалось определить чат")
		return true
	}
	*chatId = _chatId
	return
}

func (n *Node) GetChatByInvite(code string, chatId *int) (fail bool) {
	_id, err := n.repos.Chats.ChatIDByInvite(code)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*chatId = _id
	return
}

func (n *Node) MembersLimit(chatId int) (fail bool) {
	count, err := n.repos.Chats.CountMembers(chatId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxMembersOnChat {
		n.Err = resp.Error(resp.ErrBadRequest, "достигнут лимит количества участников в чате")
		return true
	}
	return
}

func (n *Node) ChatIsNotPrivate(chatId int) (fail bool) {
	if n.repos.Chats.ChatIsPrivate(chatId) {
		n.Err = resp.Error(resp.ErrBadRequest, "этот чат запривачин")
		return true
	}
	return
}

func (n *Node) UserExistsByInput(input model.LoginInput) (fail bool) {
	if !n.repos.Users.UserExistsByInput(&models.UserInput{
		Email:    input.Email,
		Password: input.Password,
	}) {
		n.Err = resp.Error(resp.ErrBadRequest, "пользователь с такими данными не найден")
		return true
	}
	return
}

func (n *Node) GetUserIDByInput(input model.LoginInput, userId *int) (fail bool) {
	_uid, err := n.repos.Users.GetUserIdByInput(&models.UserInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*userId = _uid
	return
}

func (n *Node) IsNotBanned(userId, chatId int) (fail bool) {
	if n.repos.Chats.UserIsBanned(userId, chatId) {
		n.Err = resp.Error(resp.ErrBadRequest, "вы забанены в этом чате")
		return true
	}
	return
}

func (n *Node) GetChatIDByRoom(roomId int, chatId *int) (fail bool) {
	_chatId, err := n.repos.Rooms.GetChatIDByRoomID(roomId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "не удалось найти комнату")
		return true
	}
	*chatId = _chatId
	return
}
func (n *Node) GetChatIDByAllow(allowID int, chatId *int) (fail bool) {
	_chatId, err := n.repos.Rooms.GetChatIDByAllowID(allowID)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "такого разрешения не существует")
		return true
	}
	*chatId = _chatId
	return
}

// todo *delete message*
func (n *Node) MessageAvailable(msgId, roomId int) (fail bool) {
	if !n.repos.Messages.MessageAvailableOnRoom(msgId, roomId) {
		n.Err = resp.Error(resp.ErrBadRequest, "сообщение не найдено")
		return true
	}
	return
}

func (n *Node) IsAllowedTo(action model.ActionType, roomId int, holder *models.AllowHolder) (fail bool) {
	if !n.repos.Rooms.Allowed(action, roomId, holder) {
		n.Err = resp.Error(resp.ErrBadRequest, "недостаточно прав на это действие")
		return true
	}
	return
}

func (n *Node) GetAllowHolder(userId, chatId int, holder *models.AllowHolder) (fail bool) {
	fmt.Printf("userID: %d\nchatID: %d\n", userId, chatId) // debug
	_holder, err := n.repos.Rooms.AllowHolder(userId, chatId)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "не удалось связать пользователя с чатом")
		return true
	}
	*holder = *_holder
	return
}

// deprecated
func (n *Node) IsAllowsSet(roomId int) (fail bool) {
	if !n.repos.Rooms.AllowsIsSet(roomId) {
		n.Err = resp.Error(resp.ErrBadRequest, "в комнате не установлены ограничения")
		return true
	}
	return
}

func (n *Node) GetMessageByID(msgId int, message *model.Message) (fail bool) {
	_message, err := n.repos.Messages.Message(msgId)
	if err != nil {
		n.Err = resp.Error(resp.ErrBadRequest, "сообщение не найдено")
		return true
	}
	message = _message
	return
}

func (n *Node) GetChatIDByMember(memberId int, chatId *int) (fail bool) {
	_chatId, err := n.repos.Chats.ChatIDByMemberID(memberId)
	if err != nil {
		n.Err = resp.Error(resp.ErrBadRequest, "участник не найден")
		return true
	}
	*chatId = _chatId
	return
}

func (n *Node) GetMemberBy(userId, chatId int, memberId *int) (fail bool) {
	by := n.repos.Chats.FindMemberBy(userId, chatId)
	if by == nil || *by == 0 {
		n.Err = resp.Error(resp.ErrBadRequest, "не удалось определить участника чата")
		return true
	}
	*memberId = *by
	return
}

func (n *Node) IsNotMuted(memberId int) (fail bool) {
	if n.repos.Chats.MemberIsMuted(memberId) {
		n.Err = resp.Error(resp.ErrBadRequest, "участник чата заглушен")
		return true
	}
	return
}

func (n *Node) HandleChoice(choiceString string, roomId int, handledChoice *string) (fail bool) {
	var userChoice model.UserChoice
	err := json.Unmarshal([]byte(choiceString), &userChoice)
	if err != nil {
		n.Err = resp.Error(resp.ErrBadRequest, "невалидное тело запроса")
		return true
	}
	form := n.repos.Rooms.RoomForm(roomId)

	fmt.Printf("%s", form.Fields[0].Key) // debug
	choice, aerr := matchMessageType(&userChoice, form)
	if aerr != nil {
		n.Err = resp.Error(resp.ErrBadRequest, aerr.Message)
		return true
	}
	choiceBody, err := json.Marshal(choice)
	if err != nil {
		n.Err = resp.Error(resp.ErrInternalServerError, "ошибка при обработке запроса")
		return true
	}
	*handledChoice = string(choiceBody)

	return
}

func matchMessageType(input *model.UserChoice, sample *model.Form) (*model.UserChoice, *rules.AdvancedError) {
	completed := make(map[string]string)
	for _, field := range sample.Fields {
		for _, choice := range input.Choice {
			if choice.Key == field.Key {
				var advErr *rules.AdvancedError
				if field.Length != nil && len(choice.Value) > *field.Length {
					advErr = rules.ErrChoiceValueLength
				}
				switch field.Type {
				case model.FieldTypeText:
					// nothing

				case model.FieldTypeDate:
					if _, err := strconv.ParseInt(choice.Value, 10, 64); err != nil {
						advErr = rules.ErrInvalidChoiceDate
					}

				case model.FieldTypeEmail:
					if !validator.ValidateEmail(choice.Value) {
						advErr = rules.ErrInvalidEmail
					}

				case model.FieldTypeLink:
					if !validator.ValidateLink(choice.Value) {
						advErr = rules.ErrInvalidLink
					}

				case model.FieldTypeNumeric:
					if _, err := strconv.Atoi(choice.Value); err != nil {
						advErr = rules.ErrInvalidChoiceValue
					}

				default:
					advErr = rules.ErrDataRetrieved
				}
				if advErr != nil {
					return nil, advErr
				}
				if len(field.Items) != 0 {
					contains := func(arr []string, str string) bool {
						for _, a := range arr {
							if a == str {
								return true
							}
						}
						return false
					}(field.Items, choice.Value)

					if !contains {
						return nil, rules.ErrInvalidChoiceValue
					}
				}
				completed[field.Key] = choice.Value
			}

		}
		_, ok := completed[field.Key]
		if !(ok || field.Optional) {
			return nil, rules.ErrMissingChoicePair
		}

	}
	form := &model.UserChoice{
		Choice: []*model.Case{},
	}
	for k, v := range completed {
		form.Choice = append(form.Choice, &model.Case{
			Key:   k,
			Value: v,
		})
	}
	return form, nil
}

func (n *Node) GetDefMember(memberId int, defMember *models.DefMember) (fail bool) {
	_defMember, err := n.repos.Chats.DefMember(memberId)
	if err != nil {
		n.Err = resp.Error(resp.ErrBadRequest, "участник не найден")
		return true
	}
	*defMember = _defMember
	return
}

func (n *Node) AllowNotExists(roomID int, inp *model.AllowInput) (fail bool) {
	if n.repos.Rooms.AllowExists(roomID, inp) {
		n.Err = resp.Error(resp.ErrBadRequest, "такое разрешение уже существует")
		return true
	}
	return
}
func (n *Node) ValidAllowInput(chatID int, inp *model.AllowInput) (fail bool) {
	val := inp.Value
	intVal, err := strconv.Atoi(val)
	switch inp.Group {
	case model.AllowGroupChar:
		if model.CharTypeModer.String() != val &&
			model.CharTypeAdmin.String() != val {
			n.Err = resp.Error(resp.ErrBadRequest, "невалидное значение")
			return true
		}

	case model.AllowGroupMember:
		if err != nil {
			n.Err = resp.Error(resp.ErrBadRequest, "невалидное значение")
			return true
		}
		_chatID, err := n.repos.Chats.ChatIDByMemberID(intVal)
		if err != nil || _chatID != chatID {
			println("ValidAllowInput:", err.Error())
			n.Err = resp.Error(resp.ErrBadRequest, "не удалось определить участника чата")
			return true
		}

	case model.AllowGroupRole:
		if err != nil {
			n.Err = resp.Error(resp.ErrBadRequest, "невалидное значение")
			return true
		}
		_chatID, err := n.repos.Chats.ChatIDByRoleID(intVal)
		if err != nil || _chatID != chatID {
			println("ValidAllowInput:", err.Error())
			n.Err = resp.Error(resp.ErrBadRequest, "не удалось определить роль")
			return true
		}

	default:
		println("Not implemented")
	}
	return
}

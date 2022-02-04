package piper

import (
	"encoding/json"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/validator"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

func (n Node) ChatExists(chatId int) (fail bool) {
	n.SwitchMethod("ChatExists", &bson.M{
		"chatId": chatId,
	})
	defer n.MethodTiming()

	if !n.repos.Units.UnitExistsByID(chatId, res.Chat) {
		n.SetError(resp.ErrBadRequest, "такого чата не существует")
		return true
	}
	return
}

func (n Node) UserExists(userId int) (fail bool) {
	n.SwitchMethod("UserExists", &bson.M{
		"userId": userId,
	})
	defer n.MethodTiming()

	if !n.repos.Units.UnitExistsByID(userId, res.User) {
		n.SetError(resp.ErrBadRequest, "пользователь не найден")
		return true
	}
	return
}

// ValidParams
//  with side effect
func (n Node) ValidParams(params **model.Params) (fail bool) {
	n.SwitchMethod("ValidParams", &bson.M{
		"params": params,
	})
	defer n.MethodTiming()

	if *params == nil {
		*params = &model.Params{
			Limit:  kit.IntPtr(rules.MaxLimit), // ! unsafe
			Offset: kit.IntPtr(0),
		}
		return
	}
	if (*params).Limit != nil {
		if !validator.ValidateLimit(*(*params).Limit) {
			n.SetError(resp.ErrBadRequest, "невалидное значение лимита")
			return true
		}
	} else {
		(*params).Limit = kit.IntPtr(rules.MaxLimit)
	}
	if (*params).Offset != nil {
		if !validator.ValidateOffset(*(*params).Offset) {
			n.SetError(resp.ErrBadRequest, "невалидное значение смещения")
			return true
		}
	} else {
		(*params).Offset = kit.IntPtr(0)
	}

	return
}

func (n Node) ValidNameFragment(fragment string) (fail bool) {
	n.SwitchMethod("ValidNameFragment", &bson.M{
		"fragment": fragment,
	})
	defer n.MethodTiming()

	if !validator.ValidateNameFragment(fragment) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для фрагмента имени")
		return true
	}
	return
}

func (n *Node) ValidNote(note string) (fail bool) {
	n.SwitchMethod("ValidNote", &bson.M{
		"note": note,
	})
	defer n.MethodTiming()

	if !validator.ValidateNote(note) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для заметки")
		return true
	}
	return
}

func (n Node) ValidID(id int) (fail bool) {
	n.SwitchMethod("ValidID", &bson.M{
		"id": id,
	})
	defer n.MethodTiming()

	if !validator.ValidateID(id) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для id")
		return true
	}
	return
}

func (n Node) ValidParentRoomID(id, parent int) (fail bool) {
	n.SwitchMethod("ValidParentRoomID", &bson.M{
		"id":     id,
		"parent": parent,
	})
	defer n.MethodTiming()

	if !validator.ValidateID(parent) || id == parent {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для id")
		return true
	}
	return
}

func (n Node) ValidRoomAllows(chatID int, allows *model.AllowsInput) (fail bool) {
	n.SwitchMethod("ValidRoomAllows", &bson.M{
		"allows": allows,
	})
	defer n.MethodTiming()

	if len(allows.Allows) == 0 ||
		!validator.ValidateAllowsInput(allows) ||
		!n.repos.Chats.ValidAllows(chatID, allows) {
		n.SetError(resp.ErrBadRequest, "одно из разрешений содержит ошибку")
		return true
	}
	return
}

func (n Node) ValidFindMessagesInRoom(find *model.FindMessagesInRoom) (fail bool) {
	n.SwitchMethod("ValidFindMessagesInRoom", &bson.M{
		"find": find,
	})
	defer n.MethodTiming()

	if find.Count <= 0 ||
		find.Count > rules.MaxMessagesCount ||
		find.Created == model.MessagesCreatedBefore && find.StartMessageID-find.Count < 0 {
		n.SetError(resp.ErrBadRequest, "неверное значение количества сообщений")
		return true
	}
	if !validator.ValidateID(find.StartMessageID) {
		n.SetError(resp.ErrBadRequest, "неверный ID сообщения")
		return true
	}

	return
}

func (n Node) ValidPassword(password string) (fail bool) {
	n.SwitchMethod("ValidPassword", &bson.M{
		"password": password,
	})
	defer n.MethodTiming()

	if !validator.ValidatePassword(password) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение для пароля")
		return true
	}
	return
}

func (n Node) ValidForm(form *model.UpdateFormInput) (fail bool) {
	n.SwitchMethod("ValidForm", &bson.M{
		"form": form,
	})
	defer n.MethodTiming()

	_, err := validator.ValidateRoomForm(form)
	if err != nil {
		n.SetError(resp.ErrBadRequest, err.Error())
		return true
	}
	return
}

func (n Node) OwnedLimit(userId int) (fail bool) {
	n.SwitchMethod("OwnedLimit", &bson.M{
		"userId": userId,
	})
	defer n.MethodTiming()

	count, err := n.repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxCountOwnedChats {
		n.SetError(resp.ErrBadRequest, "достигнут лимит созднных чатов")
		return true
	}
	return
}

func (n Node) ChatsLimit(userId int) (fail bool) {
	n.SwitchMethod("ChatsLimit", &bson.M{
		"userId": userId,
	})
	defer n.MethodTiming()

	count, err := n.repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxUserChats {
		n.SetError(resp.ErrBadRequest, "достигнут лимит количества чатов в которых пользователь может состоять")
		return true
	}
	return
}

func (n Node) ValidDomain(domain string) (fail bool) {
	n.SwitchMethod("ValidDomain", &bson.M{
		"domain": domain,
	})
	defer n.MethodTiming()

	if !validator.ValidateDomain(domain) {
		n.SetError(resp.ErrBadRequest, "невалидный домен")
		return true
	}
	return
}

func (n Node) ValidName(name string) (fail bool) {
	n.SwitchMethod("ValidName", &bson.M{
		"name": name,
	})
	defer n.MethodTiming()

	if !validator.ValidateName(name) {
		n.SetError(resp.ErrBadRequest, "невалидное имя")
		return true
	}
	return
}

func (n Node) DomainIsFree(domain string) (fail bool) {
	n.SwitchMethod("DomainIsFree", &bson.M{
		"domain": domain,
	})
	defer n.MethodTiming()

	if !n.repos.Units.DomainIsFree(domain) {
		n.SetError(resp.ErrBadRequest, "домен занят")
		return true
	}
	return
}

func (n Node) EmailIsFree(email string) (fail bool) {
	n.SwitchMethod("EmailIsFree", &bson.M{
		"email": email,
	})
	defer n.MethodTiming()

	if !n.repos.Users.EmailIsFree(email) {
		n.SetError(resp.ErrBadRequest, "такая почта уже занята кем-то")
		return true
	}
	return
}

func (n Node) InvitesLimit(chatId int) (fail bool) {
	n.SwitchMethod("InvitesLimit", &bson.M{
		"chatId": chatId,
	})
	defer n.MethodTiming()

	count, err := n.repos.Chats.GetCountLinks(chatId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxInviteLinks {
		n.SetError(resp.ErrBadRequest, "достигнут лимит количества инвайтов")
		return true
	}
	return
}

func (n Node) ValidInviteInput(inp model.CreateInviteInput) (fail bool) {
	n.SwitchMethod("ValidInviteInput", &bson.M{
		"inp": inp,
	})
	defer n.MethodTiming()

	if inp.Duration != nil && !validator.ValidateLifetime(*inp.Duration) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение времени жизни ссылки")
		return true
	}
	if inp.Aliens != nil && !validator.ValidateAliens(*inp.Aliens) {
		n.SetError(resp.ErrBadRequest, "недопустимое значение количества ипользований ссылки")
		return true
	}
	return
}

// IsMember does not need if the Can method is used..
func (n Node) IsMember(userId, chatId int) (fail bool) {
	n.SwitchMethod("IsMember", &bson.M{
		"userId": userId,
		"chatId": chatId,
	})
	defer n.MethodTiming()

	if !n.repos.Chats.UserIsChatMember(userId, chatId) {
		n.SetError(resp.ErrBadRequest, "пользователь не является участником чата")
		return true
	}

	return
}

func (n Node) IsNotMember(userId, chatId int) (fail bool) {
	n.SwitchMethod("IsNotMember", &bson.M{
		"userId": userId,
		"chatId": chatId,
	})
	defer n.MethodTiming()

	if n.repos.Chats.UserIsChatMember(userId, chatId) {
		n.SetError(resp.ErrBadRequest, "пользователь является участником чата")
		return true
	}
	return
}

func (n Node) RolesLimit(chatId int) (fail bool) {
	n.SwitchMethod("RolesLimit", &bson.M{
		"chatId": chatId,
	})
	defer n.MethodTiming()

	count, err := n.repos.Chats.GetCountChatRoles(chatId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxRolesInChat {
		n.SetError(resp.ErrBadRequest, "достигнут лимит количества ролей в чате")
		return true
	}
	return
}

func (n Node) RoomsLimit(chatId int) (fail bool) {
	n.SwitchMethod("RoomsLimit", &bson.M{
		"chatId": chatId,
	})
	defer n.MethodTiming()

	count, err := n.repos.Chats.GetCountRooms(chatId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxCountRooms {
		n.SetError(resp.ErrBadRequest, "достигнут лимит количества комнат в чате")
		return true
	}
	return
}

func (n Node) RoomExists(roomId int) (fail bool) {
	n.SwitchMethod("RoomExists", &bson.M{
		"roomId": roomId,
	})
	defer n.MethodTiming()

	if !n.repos.Rooms.RoomExistsByID(roomId) {
		n.SetError(resp.ErrBadRequest, "такой комнаты не существует")
		return true
	}
	return
}

func (n Node) IsNotChild(roomId int) (fail bool) {
	n.SwitchMethod("IsNotChild", &bson.M{
		"roomId": roomId,
	})
	defer n.MethodTiming()

	if n.repos.Rooms.HasParent(roomId) {
		n.SetError(resp.ErrBadRequest, "комната является веткой")
		return true
	}
	return
}

func (n Node) HasInvite(chatId int, code string) (fail bool) {
	n.SwitchMethod("HasInvite", &bson.M{
		"chatId": chatId,
		"code":   code,
	})
	defer n.MethodTiming()

	if !n.repos.Chats.HasInvite(chatId, code) {
		n.SetError(resp.ErrBadRequest, "такого кода не существует")
		return true
	}
	return
}

func (n Node) InviteIsRelevant(code string) (fail bool) {
	n.SwitchMethod("InviteIsRelevant", &bson.M{
		"code": code,
	})
	defer n.MethodTiming()

	if !n.repos.Chats.InviteIsRelevant(code) {
		n.SetError(resp.ErrBadRequest, "инвайт неактуален")
		return true
	}
	return
}

func (n Node) RoleExists(chatID, roleID int) (fail bool) {
	n.SwitchMethod("RoleExists", &bson.M{
		"chatID": chatID,
		"roleID": roleID,
	})
	defer n.MethodTiming()

	if !n.repos.Chats.RoleExistsByID(chatID, roleID) {
		n.SetError(resp.ErrBadRequest, "такой роли не существует")
		return true
	}
	return
}
func (n Node) InviteIsExists(code string) (fail bool) {
	n.SwitchMethod("InviteIsExists", &bson.M{
		"code": code,
	})
	defer n.MethodTiming()

	if !n.repos.Chats.InviteExistsByCode(code) {
		n.SetError(resp.ErrBadRequest, "инвайта не существует")
		return true
	}
	return
}
func (n Node) GetChatIDByRole(roleId int, chatId *int) (fail bool) {
	n.SwitchMethod("GetChatIDByRole", &bson.M{
		"roleId": roleId,
		"chatId": chatId,
	})
	defer n.MethodTiming()

	_chatId, err := n.repos.Chats.ChatIDByRoleID(roleId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "не удалось определить чат")
		return true
	}
	*chatId = _chatId
	return
}

func (n Node) GetChatByInvite(code string, chatId *int) (fail bool) {
	n.SwitchMethod("GetChatByInvite", &bson.M{
		"code":   code,
		"chatId": chatId,
	})
	defer n.MethodTiming()

	_id, err := n.repos.Chats.ChatIDByInvite(code)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*chatId = _id
	return
}

func (n Node) MembersLimit(chatId int) (fail bool) {
	n.SwitchMethod("MembersLimit", &bson.M{
		"chatId": chatId,
	})
	defer n.MethodTiming()

	count, err := n.repos.Chats.CountMembers(chatId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	if count >= rules.MaxMembersOnChat {
		n.SetError(resp.ErrBadRequest, "достигнут лимит количества участников в чате")
		return true
	}
	return
}

func (n Node) ChatIsNotPrivate(chatId int) (fail bool) {
	n.SwitchMethod("ChatIsNotPrivate", &bson.M{
		"chatId": chatId,
	})
	defer n.MethodTiming()

	if n.repos.Chats.ChatIsPrivate(chatId) {
		n.SetError(resp.ErrBadRequest, "этот чат запривачин")
		return true
	}
	return
}

func (n Node) UserExistsByRequisites(input *models.LoginRequisites) (fail bool) {
	n.SwitchMethod("UserExistsByRequisites", &bson.M{
		"input": input,
	})
	defer n.MethodTiming()

	if !n.repos.Users.UserExistsByRequisites(input) {
		n.SetError(resp.ErrBadRequest, "нееверный логин или пароль ")
		return true
	}
	return
}

func (n Node) GetUserIDByRequisites(input *models.LoginRequisites, userId *int) (fail bool) {
	n.SwitchMethod("GetUserIDByRequisites", &bson.M{
		"input":  input,
		"userId": userId,
	})
	defer n.MethodTiming()

	_uid, err := n.repos.Users.GetUserIdByRequisites(input)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка базы данных")
		return true
	}
	*userId = _uid
	return
}

func (n Node) IsNotBanned(userId, chatId int) (fail bool) {
	n.SwitchMethod("IsNotBanned", &bson.M{
		"userId": userId,
		"chatId": chatId,
	})
	defer n.MethodTiming()

	if n.repos.Chats.UserIsBanned(userId, chatId) {
		n.SetError(resp.ErrBadRequest, "вы забанены в этом чате")
		return true
	}
	return
}

func (n Node) GetChatIDByRoom(roomId int, chatId *int) (fail bool) {
	n.SwitchMethod("GetChatIDByRoom", &bson.M{
		"roomId": roomId,
		"chatId": chatId,
	})
	defer n.MethodTiming()

	_chatId, err := n.repos.Rooms.GetChatIDByRoomID(roomId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "не удалось найти комнату")
		return true
	}
	*chatId = _chatId
	return
}
func (n Node) GetChatIDByAllow(allowID int, chatId *int) (fail bool) {
	n.SwitchMethod("GetChatIDByAllow", &bson.M{
		"allowID": allowID,
		"chatId":  chatId,
	})
	defer n.MethodTiming()

	_chatId, err := n.repos.Rooms.GetChatIDByAllowID(allowID)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "такого разрешения не существует")
		return true
	}
	*chatId = _chatId
	return
}

// todo *delete message*
func (n Node) MessageAvailable(msgId, roomId int) (fail bool) {
	n.SwitchMethod("MessageAvailable", &bson.M{
		"msgId":  msgId,
		"roomId": roomId,
	})
	defer n.MethodTiming()

	if !n.repos.Messages.MessageAvailableOnRoom(msgId, roomId) {
		n.SetError(resp.ErrBadRequest, "сообщение не найдено")
		return true
	}
	return
}

func (n Node) IsAllowedTo(action model.ActionType, roomId int, holder *models.AllowHolder) (fail bool) {
	n.SwitchMethod("IsAllowedTo", &bson.M{
		"action": action,
		"roomId": roomId,
		"holder": holder,
	})
	defer n.MethodTiming()

	if !n.repos.Rooms.Allowed(action, roomId, holder) {
		n.SetError(resp.ErrBadRequest, "недостаточно прав на это действие")
		return true
	}
	return
}

func (n Node) GetAllowHolder(userId, chatId int, holder *models.AllowHolder) (fail bool) {
	n.SwitchMethod("GetAllowHolder", &bson.M{
		"userId": userId,
		"chatId": chatId,
		"holder": holder,
	})
	defer n.MethodTiming()

	fmt.Printf("userID: %d\nchatID: %d\n", userId, chatId) // debug
	_holder, err := n.repos.Rooms.AllowHolder(userId, chatId)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "не удалось связать пользователя с чатом")
		return true
	}
	*holder = *_holder
	return
}

// deprecated
func (n Node) IsAllowsSet(roomId int) (fail bool) {
	n.SwitchMethod("IsAllowsSet", &bson.M{
		"roomId": roomId,
	})
	defer n.MethodTiming()

	if !n.repos.Rooms.AllowsIsSet(roomId) {
		n.SetError(resp.ErrBadRequest, "в комнате не установлены ограничения")
		return true
	}
	return
}

func (n *Node) GetMessageByID(msgId int, message *model.Message) (fail bool) {
	n.SwitchMethod("GetMessageByID", &bson.M{
		"msgId":   msgId,
		"message": message,
	})
	defer n.MethodTiming()

	_message, err := n.repos.Messages.Message(msgId)
	if err != nil {
		n.SetError(resp.ErrBadRequest, "сообщение не найдено")
		return true
	}
	message = _message
	return
}

func (n Node) GetChatIDByMember(memberId int, chatId *int) (fail bool) {
	n.SwitchMethod("GetChatIDByMember", &bson.M{
		"memberId": memberId,
		"chatId":   chatId,
	})
	defer n.MethodTiming()

	_chatId, err := n.repos.Chats.ChatIDByMemberID(memberId)
	if err != nil {
		n.SetError(resp.ErrBadRequest, "участник не найден")
		return true
	}
	*chatId = _chatId
	return
}

func (n Node) GetMemberBy(userId, chatId int, memberId *int) (fail bool) {
	n.SwitchMethod("GetMemberBy", &bson.M{
		"userId":   userId,
		"chatId":   chatId,
		"memberId": memberId,
	})
	defer n.MethodTiming()

	by := n.repos.Chats.FindMemberBy(userId, chatId)
	if by == nil || *by == 0 {
		n.SetError(resp.ErrBadRequest, "мембер не найден")
		return true
	}
	*memberId = *by
	return
}

func (n Node) IsNotMuted(memberId int) (fail bool) {
	n.SwitchMethod("IsNotMuted", &bson.M{
		"memberId": memberId,
	})
	defer n.MethodTiming()

	if n.repos.Chats.MemberIsMuted(memberId) {
		n.SetError(resp.ErrBadRequest, "участник заглушен")
		return true
	}
	return
}

func (n Node) HandleChoice(choiceString string, roomId int, handledChoice *string) (fail bool) {
	n.SwitchMethod("HandleChoice", &bson.M{
		"choiceString":  choiceString,
		"roomId":        roomId,
		"handledChoice": handledChoice,
	})
	defer n.MethodTiming()

	var userChoice model.UserChoice
	err := json.Unmarshal([]byte(choiceString), &userChoice)
	if err != nil {
		n.SetError(resp.ErrBadRequest, "невалидное тело запроса")
		return true
	}
	form := n.repos.Rooms.RoomForm(roomId)

	fmt.Printf("%s", form.Fields[0].Key) // debug
	choice, aerr := matchMessageType(&userChoice, form)
	if aerr != nil {
		n.SetError(resp.ErrBadRequest, aerr.Message)
		return true
	}
	choiceBody, err := json.Marshal(choice)
	if err != nil {
		n.SetError(resp.ErrInternalServerError, "ошибка при обработке запроса")
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

func (n Node) GetDefMember(memberId int, defMember *models.DefMember) (fail bool) {
	n.SwitchMethod("GetDefMember", &bson.M{
		"memberId":  memberId,
		"defMember": defMember,
	})
	defer n.MethodTiming()

	_defMember, err := n.repos.Chats.DefMember(memberId)
	if err != nil {
		n.SetError(resp.ErrBadRequest, "участник не найден")
		return true
	}
	*defMember = _defMember
	return
}

func (n Node) AllowsNotExists(roomID int, inp *model.AllowsInput) (fail bool) {
	n.SwitchMethod("AllowsNotExists", &bson.M{
		"roomID": roomID,
		"inp":    inp,
	})
	defer n.MethodTiming()

	if n.repos.Rooms.AllowsExists(roomID, inp) {
		n.SetError(resp.ErrBadRequest, "такое разрешение уже существует")
		return true
	}
	return
}

func (n Node) ValidAllowInput(chatID int, inp *model.AllowInput) (fail bool) {
	n.SwitchMethod("ValidAllowInput", &bson.M{
		"chatID": chatID,
		"inp":    inp,
	})
	defer n.MethodTiming()

	val := inp.Value
	intVal, err := strconv.Atoi(val)
	switch inp.Group {
	case model.AllowGroupChar:
		if model.CharTypeModer.String() != val &&
			model.CharTypeAdmin.String() != val {
			n.SetError(resp.ErrBadRequest, "невалидное значение")
			return true
		}

	case model.AllowGroupMember:
		if err != nil {
			n.SetError(resp.ErrBadRequest, "невалидное значение")
			return true
		}
		_chatID, err := n.repos.Chats.ChatIDByMemberID(intVal)
		if err != nil || _chatID != chatID {
			println("ValidAllowInput:", err.Error())
			n.SetError(resp.ErrBadRequest, "не удалось определить участника чата")
			return true
		}

	case model.AllowGroupRole:
		if err != nil {
			n.SetError(resp.ErrBadRequest, "невалидное значение")
			return true
		}
		_chatID, err := n.repos.Chats.ChatIDByRoleID(intVal)
		if err != nil || _chatID != chatID {
			println("ValidAllowInput:", err.Error())
			n.SetError(resp.ErrBadRequest, "не удалось определить роль")
			return true
		}

	default:
		println("Not implemented")
	}
	return
}

func (n Node) UserHasAccessToChats(userID int, chats *[]int, submembers **[]*models.SubUser) (fail bool) {
	n.SwitchMethod("UserHasAccessToChats", &bson.M{
		"userID":     userID,
		"chats":      chats,
		"submembers": submembers,
	})
	defer n.MethodTiming()

	if !validator.ValidateIDs(*chats) {
		n.SetError(resp.ErrBadRequest, "невалидный id")
		return true
	}
	members, yes := n.repos.Chats.UserHasAccessToChats(userID, chats)
	if !yes {
		n.SetError(resp.ErrBadRequest, "нет доступа к одному из чатов")
		return true
	}
	*submembers = &members
	return
}

func (n Node) ValidSessionKey(sessionKey string) (fail bool) {
	n.SwitchMethod("ValidSessionKey", &bson.M{
		"sessionKey": sessionKey,
	})
	defer n.MethodTiming()

	if !validator.ValidateSessionKey(sessionKey) {
		n.SetError(resp.ErrBadRequest, "не валидный ключ сессии")
		return true
	}
	return
}

func (n Node) ValidEmail(email string) (fail bool) {
	n.SwitchMethod("ValidateEmail", &bson.M{
		"email": email,
	})
	defer n.MethodTiming()

	if !validator.ValidateEmail(email) {
		n.SetError(resp.ErrBadRequest, "не валидный email")
		return true
	}
	return
}

func (n Node) GetRegistrationSession(email, code string, regi **models.RegisterData) (fail bool) {
	n.SwitchMethod("GetRegistrationSession", &bson.M{
		"email": email,
		"code":  code,
		"regi":  regi,
	})
	defer n.MethodTiming()

	var err error
	*regi, err = n.repos.Users.GetRegistrationSession(email, code)
	if err != nil {
		n.SetError(resp.ErrBadRequest, "сессии не существует")
		return true
	}
	return
}

func (n Node) ValidRegisterInput(input *model.RegisterInput) (fail bool) {
	n.SwitchMethod("ValidRegisterInput", &bson.M{
		"input": input,
	})
	defer n.MethodTiming()

	switch {
	case !validator.ValidateDomain(input.Domain):
		n.SetError(resp.ErrBadRequest, "домен не соответствует требованиям")
		return true

	case !validator.ValidateName(input.Name):
		n.SetError(resp.ErrBadRequest, "имя не соответствует требованиям")
		return true

	case !validator.ValidateEmail(input.Email):
		n.SetError(resp.ErrBadRequest, "имеил не соответствует требованиям")
		return true

	case !validator.ValidatePassword(input.Password):
		n.SetError(resp.ErrBadRequest, "пароль не соответствует требованиям")
		return true

	}

	return
}

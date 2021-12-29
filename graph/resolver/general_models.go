package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/tlog"
)

func (r *chatResolver) Owner(ctx context.Context, obj *model.Chat) (model.UserResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	chatID := obj.Unit.ID
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveOwner(clientID, chatID) {
		return pl.Err, nil
	}

	user, err := r.Services.Repos.Chats.Owner(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return user, nil
}

func (r *chatResolver) Rooms(ctx context.Context, obj *model.Chat) (model.RoomsResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveRooms(clientID, chatID) {
		return pl.Err, nil
	}

	rooms, err := r.Services.Repos.Chats.Rooms(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	for _, room := range rooms.Rooms {
		room.Chat = obj
	}

	return rooms, nil
}

func (r *chatResolver) CountMembers(ctx context.Context, obj *model.Chat) (model.CountMembersResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveCountMembers(clientID, chatID) {
		return pl.Err, nil
	}
	count, err := r.Services.Repos.Chats.CountMembers(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return model.IntValue{Value: &count}, nil
}

func (r *chatResolver) Members(ctx context.Context, obj *model.Chat) (model.MembersResult, error) {
	chatID := obj.Unit.ID
	tl := tlog.Start("chatResolver > Members [cid:" + strconv.Itoa(chatID) + "]")
	defer tl.Fine()
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveMembers(clientID, chatID) {
		return pl.Err, nil
	}

	members, err := r.Services.Repos.Chats.Members(chatID)
	if err != nil {
		println(err.Error()) // debug
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	return members, nil
}

func (r *chatResolver) Roles(ctx context.Context, obj *model.Chat) (model.RolesResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveRoles(clientID, chatID) {
		return pl.Err, nil
	}
	_roles, err := r.Services.Repos.Chats.Roles(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	var roles model.Roles
	for _, role := range _roles {
		roles.Roles = append(roles.Roles, &model.Role{
			ID:    role.ID,
			Name:  role.Name,
			Color: role.Color,
		})
	}
	return roles, nil
}

func (r *chatResolver) Invites(ctx context.Context, obj *model.Chat) (model.InvitesResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveInvites(clientID, chatID) {
		return pl.Err, nil
	}

	_invites, err := r.Services.Repos.Chats.Invites(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	var invites model.Invites
	for _, i := range _invites.Invites {
		invites.Invites = append(invites.Invites, &model.Invite{
			Code:      i.Code,
			Aliens:    &i.Aliens,
			ExpiresAt: &i.Exp,
		})
	}
	return invites, nil
}

func (r *chatResolver) Banlist(ctx context.Context, obj *model.Chat) (model.UsersResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveBanlist(clientID, chatID) {
		return pl.Err, nil
	}

	banlist, err := r.Services.Repos.Chats.Banlist(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	var users model.Users
	for _, user := range banlist.Users {
		users.Users = append(users.Users, &model.User{
			Unit: &model.Unit{
				ID:     user.Unit.ID,
				Domain: user.Unit.Domain,
				Name:   user.Unit.Name,
				Type:   model.UnitType(user.Unit.Type),
			},
		})
	}
	return users, nil
}

func (r *chatResolver) Me(ctx context.Context, obj *model.Chat) (model.MemberResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.IsMember(clientID, chatID) {
		return pl.Err, nil
	}

	member, err := r.Services.Repos.Chats.MemberBy(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	member.Chat = obj
	return member, nil
}

func (r *inviteInfoResolver) CountMembers(ctx context.Context, obj *model.InviteInfo) (model.CountMembersResult, error) {
	chatID := obj.Unit.ID

	count, err := r.Services.Repos.Chats.CountMembers(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return model.IntValue{Value: &count}, nil
}

func (r *meResolver) Chats(ctx context.Context, obj *model.Me) ([]*model.Chat, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	_chats, err := r.Services.Repos.Users.Chats(clientID)
	if err != nil {
		return nil, nil // resp.Error
	}
	var chats []*model.Chat
	for _, chat := range _chats {
		chats = append(chats, &model.Chat{
			Unit: &model.Unit{
				ID:     chat.Unit.ID,
				Domain: chat.Unit.Domain,
				Name:   chat.Unit.Name,
				Type:   model.UnitType(chat.Unit.Type),
			},
			Private: chat.Private,
		})
	}
	return chats, nil
}

func (r *meResolver) OwnedChats(ctx context.Context, obj *model.Me) ([]*model.Chat, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	_chats, err := r.Services.Repos.Users.OwnedChats(clientID)
	if err != nil {
		return nil, nil // resp.Error
	}
	var chats []*model.Chat
	for _, chat := range _chats {
		chats = append(chats, &model.Chat{
			Unit: &model.Unit{
				ID:     chat.Unit.ID,
				Domain: chat.Unit.Domain,
				Name:   chat.Unit.Name,
				Type:   model.UnitType(chat.Unit.Type),
			},
			Private: chat.Private,
		})
	}
	return chats, nil
}

func (r *memberResolver) Chat(ctx context.Context, obj *model.Member) (*model.Chat, error) {
	chatID := obj.Chat.Unit.ID
	chat, err := r.Services.Repos.Chats.Chat(chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (r *memberResolver) Role(ctx context.Context, obj *model.Member) (model.RoleResult, error) {
	memberID := obj.ID
	tl := tlog.Start("memberResolver > Role [mid:" + strconv.Itoa(memberID) + "]")
	defer tl.Fine()
	role := r.Services.Repos.Chats.MemberRole(memberID)

	return role, nil
}

func (r *messageResolver) Room(ctx context.Context, obj *model.Message) (*model.Room, error) {
	roomID := obj.Room.RoomID
	room, err := r.Services.Repos.Rooms.Room(roomID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *messageResolver) ReplyTo(ctx context.Context, obj *model.Message) (*model.Message, error) {
	if obj.ReplyTo == nil {
		return nil, nil
	}
	message, err := r.Services.Repos.Messages.Message(obj.ReplyTo.ID)
	if err != nil {
		return nil, err
	}
	return message, err
}

func (r *messageResolver) Author(ctx context.Context, obj *model.Message) (*model.Member, error) {
	if obj.Author == nil {
		return nil, nil
	}
	memberID := obj.Author.ID
	member, err := r.Services.Repos.Chats.Member(memberID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *roomResolver) Chat(ctx context.Context, obj *model.Room) (*model.Chat, error) {
	chatID := obj.Chat.Unit.ID
	tl := tlog.Start("roomResolver > Chat [cid:" + strconv.Itoa(chatID) + "]")
	chat, err := r.Services.Repos.Chats.Chat(chatID)
	if err != nil {
		tl.FineWithReason(err.Error())
		return nil, err
	}
	tl.Fine()
	return chat, nil
}

func (r *roomResolver) Form(ctx context.Context, obj *model.Room) (model.RoomFormResult, error) {
	roomID := obj.RoomID
	tl := tlog.Start("roomResolver > Form [rid:" + strconv.Itoa(roomID) + "]")
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	tl.TimeWithStatus("read clientID")
	node := r.Piper.CreateNode()
	//time.Sleep(10 * time.Millisecond) // чтобы ускорть создание ноды, саиме пошел на небольшую хитрость
	defer node.Kill()
	tl.TimeWithStatus("create node")
	var (
		chatID = obj.Chat.Unit.ID
		holder models.AllowHolder
	)
	tl.TimeWithStatus("definition vars")
	if node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(rules.AllowRead, roomID, &holder) {
		tl.FineWithReason(node.Err.Error)
		return node.Err, nil
	}
	tl.TimeWithStatus("pipline finally")
	form := r.Services.Repos.Rooms.RoomForm(roomID)
	tl.TimeWithStatus("final db query")
	tl.Fine()
	return form, nil
}

func (r *roomResolver) Allows(ctx context.Context, obj *model.Room) (model.AllowsResult, error) {
	allows, err := r.Services.Repos.Rooms.GetAllows(obj.RoomID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	allows.Room = obj
	for _, member := range append(allows.AllowRead.Members.Members, allows.AllowWrite.Members.Members...) {
		member.Chat = obj.Chat
	}

	return allows, nil
}

func (r *roomResolver) Messages(ctx context.Context, obj *model.Room, find model.FindMessagesInRoomByUnionInput, params model.Params) (model.MessagesResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	var (
		roomID = obj.RoomID
		chatID = obj.Chat.Unit.ID
		holder models.AllowHolder
	)

	pl := piping.NewPipeline(r.Services.Repos)
	if pl.GetAllowHolder(clientID, chatID, &holder) ||
		pl.IsAllowedTo(rules.AllowRead, roomID, &holder) {
		return pl.Err, nil
	}
	room := r.Services.Repos.Messages.MessagesFromRoom(roomID, chatID, &find, &params)

	return room, nil
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// InviteInfo returns generated.InviteInfoResolver implementation.
func (r *Resolver) InviteInfo() generated.InviteInfoResolver { return &inviteInfoResolver{r} }

// Me returns generated.MeResolver implementation.
func (r *Resolver) Me() generated.MeResolver { return &meResolver{r} }

// Member returns generated.MemberResolver implementation.
func (r *Resolver) Member() generated.MemberResolver { return &memberResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type chatResolver struct{ *Resolver }
type inviteInfoResolver struct{ *Resolver }
type meResolver struct{ *Resolver }
type memberResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type roomResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *roomResolver) MsgFormat(ctx context.Context, obj *model.Room) (*model.Form, error) {
	panic(fmt.Errorf("not implemented"))
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *chatResolver) Owner(ctx context.Context, obj *model.Chat) (model.UserResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	chatID := obj.Unit.ID
	pl := piping.NewPipeline(ctx, r.Services.Repos)
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
	pl := piping.NewPipeline(ctx, r.Services.Repos)
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
	pl := piping.NewPipeline(ctx, r.Services.Repos)
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
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.IsMember(clientID, chatID) ||
		pl.Can.ObserveMembers(clientID, chatID) {
		return pl.Err, nil
	}

	members, err := r.Services.Repos.Chats.Members(chatID)
	if err != nil {
		println(err.Error())
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	for _, member := range members.Members {
		member.Chat = obj
	}
	ctx = context.WithValue(ctx, rules.ChatIDFromChat, chatID)
	return members, nil
}

func (r *chatResolver) Roles(ctx context.Context, obj *model.Chat) (model.RolesResult, error) {
	chatID := obj.Unit.ID
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
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
	pl := piping.NewPipeline(ctx, r.Services.Repos)
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
	pl := piping.NewPipeline(ctx, r.Services.Repos)
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
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.IsMember(clientID, chatID) {
		return pl.Err, nil
	}

	member, err := r.Services.Repos.Chats.Member(clientID, chatID)
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

func (r *memberResolver) Role(ctx context.Context, obj *model.Member) (model.RoleResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	chatID := obj.Chat.Unit.ID
	userID := obj.User.Unit.ID
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if
	// pl.IsMember(clientID, chatID) ||
	//pl.HasRole(userID, chatID)  ||
	pl.Can.ObserveRoles(clientID, chatID) {
		return pl.Err, nil
	}
	role, err := r.Services.Repos.Chats.UserRole(userID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	return role, nil
}

func (r *messageResolver) Room(ctx context.Context, obj *model.Message) (*model.Room, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *messageResolver) ReplyTo(ctx context.Context, obj *model.Message) (*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *messageResolver) Author(ctx context.Context, obj *model.Message) (*model.Unit, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roomResolver) MsgFormat(ctx context.Context, obj *model.Room) (*model.Form, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roomResolver) Allows(ctx context.Context, obj *model.Room) (model.AllowsResult, error) {
	allows, err := r.Services.Repos.Rooms.GetAllows(obj.ID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	allows.Room = obj
	for _, member := range append(allows.AllowRead.Members.Members, allows.AllowWrite.Members.Members...) {
		member.Chat = obj.Chat
	}

	return allows, nil
}

func (r *roomResolver) Messages(ctx context.Context, obj *model.Room) ([]*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// InviteInfo returns generated.InviteInfoResolver implementation.
func (r *Resolver) InviteInfo() generated.InviteInfoResolver { return &inviteInfoResolver{r} }

// Member returns generated.MemberResolver implementation.
func (r *Resolver) Member() generated.MemberResolver { return &memberResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type chatResolver struct{ *Resolver }
type inviteInfoResolver struct{ *Resolver }
type memberResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type roomResolver struct{ *Resolver }

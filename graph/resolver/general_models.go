package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *chatResolver) Owner(ctx context.Context, obj *model.Chat) (model.UserResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   = obj.Unit.ID
	)

	if node.IsMember(clientID, chatID) {
		return node.Err, nil
	}

	user, err := r.Services.Repos.Chats.Owner(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return user, nil
}

func (r *chatResolver) Rooms(ctx context.Context, obj *model.Chat) (model.RoomsResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.Err, nil
	}

	rooms, err := r.Services.Repos.Chats.Rooms(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return rooms, nil
}

func (r *chatResolver) Members(ctx context.Context, obj *model.Chat) (model.MembersResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.Err, nil
	}

	members, err := r.Services.Repos.Chats.Members(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return members, nil
}

func (r *chatResolver) Roles(ctx context.Context, obj *model.Chat) (model.RolesResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.Err, nil
	}
	roles, err := r.Services.Repos.Chats.Roles(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return roles, nil
}

func (r *chatResolver) Invites(ctx context.Context, obj *model.Chat) (model.InvitesResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) ||
		node.CanObserveInvites(clientID, chatID) {
		return node.Err, nil
	}

	invites, err := r.Services.Repos.Chats.Invites(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return invites, nil
}

func (r *chatResolver) Banlist(ctx context.Context, obj *model.Chat) (model.UsersResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) ||
		node.CanObserveBanlist(clientID, chatID) {
		return node.Err, nil
	}

	users, err := r.Services.Repos.Chats.Banlist(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return users, nil
}

func (r *chatResolver) Me(ctx context.Context, obj *model.Chat) (model.MemberResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.Err, nil
	}

	member, err := r.Services.Repos.Chats.MemberBy(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return member, nil
}

func (r *meResolver) Chats(ctx context.Context, obj *model.Me) (*model.Chats, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	chats, err := r.Services.Repos.Users.Chats(clientID)
	if err != nil {
		return nil, nil // todo resp.Error
	}

	return chats, nil
}

func (r *meResolver) OwnedChats(ctx context.Context, obj *model.Me) (*model.Chats, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	chats, err := r.Services.Repos.Users.OwnedChats(clientID)
	if err != nil {
		return nil, nil // todo resp.Error
	}

	return chats, nil
}

func (r *memberResolver) Chat(ctx context.Context, obj *model.Member) (*model.Chat, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	chatID := obj.Chat.Unit.ID

	chat, err := r.Services.Repos.Chats.Chat(chatID)
	if err != nil {
		return nil, nil // todo resp.Error
	}

	return chat, nil
}

func (r *memberResolver) Role(ctx context.Context, obj *model.Member) (model.RoleResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	memberID := obj.ID

	role := r.Services.Repos.Chats.MemberRole(memberID)

	return role, nil
}

func (r *messageResolver) Room(ctx context.Context, obj *model.Message) (*model.Room, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	roomID := obj.Room.RoomID

	room, err := r.Services.Repos.Rooms.Room(roomID)
	if err != nil {
		return nil, err // todo resp.Error
	}
	return room, nil
}

func (r *messageResolver) ReplyTo(ctx context.Context, obj *model.Message) (*model.Message, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	if obj.ReplyTo == nil {
		return nil, nil // так и надо
	}

	message, err := r.Services.Repos.Messages.Message(obj.ReplyTo.ID)
	if err != nil {
		return nil, err // todo resp.Error
	}

	return message, nil
}

func (r *messageResolver) User(ctx context.Context, obj *model.Message) (*model.User, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	if obj.User == nil {
		return nil, nil // так и надо
	}

	userID := obj.User.Unit.ID

	user, err := r.Services.Repos.Users.User(userID)
	if err != nil {
		return nil, err // todo resp.Error
	}

	return user, nil
}

func (r *roomResolver) Chat(ctx context.Context, obj *model.Room) (*model.Chat, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	chatID := obj.Chat.Unit.ID

	chat, err := r.Services.Repos.Chats.Chat(chatID)
	if err != nil {
		return nil, err // todo resp.Error
	}

	return chat, nil
}

func (r *roomResolver) Form(ctx context.Context, obj *model.Room) (model.RoomFormResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		roomID   = obj.RoomID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   = obj.Chat.Unit.ID
		holder   models.AllowHolder
	)

	if node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(model.ActionTypeRead, roomID, &holder) {
		return node.Err, nil
	}
	form := r.Services.Repos.Rooms.RoomForm(roomID)

	return form, nil
}

func (r *roomResolver) Allows(ctx context.Context, obj *model.Room) (model.AllowsResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	allows, err := r.Services.Repos.Rooms.Allows(obj.RoomID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return allows, nil
}

func (r *roomResolver) Messages(ctx context.Context, obj *model.Room, find model.FindMessagesInRoom) (model.MessagesResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		roomID   = obj.RoomID
		chatID   = obj.Chat.Unit.ID
		holder   models.AllowHolder
	)

	if node.ValidFindMessagesInRoom(&find) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(model.ActionTypeRead, roomID, &holder) {
		return node.Err, nil
	}
	room := r.Services.Repos.Messages.MessagesFromRoom(roomID, chatID, &find)

	return room, nil
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// Me returns generated.MeResolver implementation.
func (r *Resolver) Me() generated.MeResolver { return &meResolver{r} }

// Member returns generated.MemberResolver implementation.
func (r *Resolver) Member() generated.MemberResolver { return &memberResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type chatResolver struct{ *Resolver }
type meResolver struct{ *Resolver }
type memberResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type roomResolver struct{ *Resolver }

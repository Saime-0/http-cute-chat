package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *chatResolver) Owner(ctx context.Context, obj *model.Chat) (model.UserResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Owner", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   = obj.Unit.ID
	)

	if node.IsMember(clientID, chatID) {
		return node.GetError(), nil
	}

	user, err := r.Services.Repos.Chats.Owner(chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return user, nil
}

func (r *chatResolver) Rooms(ctx context.Context, obj *model.Chat) (model.RoomsResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Rooms", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.GetError(), nil
	}

	rooms, err := r.Dataloader.Rooms(chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return rooms, nil
}

func (r *chatResolver) Members(ctx context.Context, obj *model.Chat) (model.MembersResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Members", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.GetError(), nil
	}

	members, err := r.Services.Repos.Chats.Members(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return members, nil
}

func (r *chatResolver) Roles(ctx context.Context, obj *model.Chat) (model.RolesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Roles", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.GetError(), nil
	}
	roles, err := r.Services.Repos.Chats.Roles(chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return roles, nil
}

func (r *chatResolver) Invites(ctx context.Context, obj *model.Chat) (model.InvitesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Invites", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) ||
		node.CanObserveInvites(clientID, chatID) {
		return node.GetError(), nil
	}

	invites, err := r.Services.Repos.Chats.Invites(chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return invites, nil
}

func (r *chatResolver) Banlist(ctx context.Context, obj *model.Chat) (model.UsersResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Banlist", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) ||
		node.CanObserveBanlist(clientID, chatID) {
		return node.GetError(), nil
	}

	users, err := r.Services.Repos.Chats.Banlist(chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return users, nil
}

func (r *chatResolver) Me(ctx context.Context, obj *model.Chat) (model.MemberResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chat.Me", &bson.M{
		"chatID (obj.Unit.ID)": obj.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		chatID   = obj.Unit.ID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.IsMember(clientID, chatID) {
		return node.GetError(), nil
	}

	member, err := r.Services.Repos.Chats.MemberBy(clientID, chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return member, nil
}

func (r *listenCollectionResolver) Collection(ctx context.Context, obj *model.ListenCollection) ([]*model.ListenedChat, error) {
	collection := r.Subix.ClientCollection(obj.SessionKey)
	return collection, nil
}

func (r *meResolver) Chats(ctx context.Context, obj *model.Me) (*model.Chats, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Me.Chats", nil)
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	chats, err := r.Services.Repos.Users.Chats(clientID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}

	return chats, nil
}

func (r *meResolver) OwnedChats(ctx context.Context, obj *model.Me) (*model.Chats, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Me.OwnedChats", nil)
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	chats, err := r.Services.Repos.Users.OwnedChats(clientID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}

	return chats, nil
}

func (r *memberResolver) Chat(ctx context.Context, obj *model.Member) (*model.Chat, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Member.Chat", &bson.M{
		"chatID (obj.Chat.Unit.ID)": obj.Chat.Unit.ID,
	})
	defer node.MethodTiming()

	chatID := obj.Chat.Unit.ID

	chat, err := r.Services.Repos.Chats.Chat(chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}

	return chat, nil
}

func (r *memberResolver) Role(ctx context.Context, obj *model.Member) (model.RoleResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Member.Role", &bson.M{
		"memberID (obj.ID)": obj.Chat.Unit.ID,
	})
	defer node.MethodTiming()

	memberID := obj.ID

	role, err := r.Dataloader.MemberRole(memberID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}

	return role, nil
}

func (r *messageResolver) Room(ctx context.Context, obj *model.Message) (*model.Room, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.Room", &bson.M{
		"roomID (obj.Room.RoomID)": obj.Room.RoomID,
	})
	defer node.MethodTiming()

	roomID := obj.Room.RoomID

	//room, err := r.Services.Repos.Rooms.Room(roomID)
	room, err := r.Dataloader.Room(roomID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}
	return room, nil
}

func (r *messageResolver) ReplyTo(ctx context.Context, obj *model.Message) (*model.Message, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.ReplyTo", &bson.M{
		"obj.ReplyTo": obj.ReplyTo,
	})
	defer node.MethodTiming()

	if obj.ReplyTo == nil {
		return nil, nil // так и надо
	}

	//message, err := r.Services.Repos.Messages.Message(obj.ReplyTo.ID)
	message, err := r.Dataloader.Message(obj.ReplyTo.ID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}

	return message, nil
}

func (r *messageResolver) User(ctx context.Context, obj *model.Message) (*model.User, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Message.User", &bson.M{
		"obj.User": obj.User,
	})
	defer node.MethodTiming()

	if obj.User == nil {
		return nil, nil // так и надо
	}

	userID := obj.User.Unit.ID

	user, err := r.Dataloader.User(userID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return nil, cerrors.New("произошла ошибка во время обработки данных") // todo resp.Error
	}

	return user, nil
}

func (r *roomResolver) Chat(ctx context.Context, obj *model.Room) (*model.Chat, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Room.Chat", &bson.M{
		"chatID (obj.Chat.Unit.ID)": obj.Chat.Unit.ID,
	})
	defer node.MethodTiming()

	chatID := obj.Chat.Unit.ID

	chat, err := r.Services.Repos.Chats.Chat(chatID)
	if err != nil {
		return nil, err // todo resp.Error
	}

	return chat, nil
}

func (r *roomResolver) Form(ctx context.Context, obj *model.Room) (model.RoomFormResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Room.Form", &bson.M{
		"roomID (obj.RoomID)":       obj.RoomID,
		"chatID (obj.Chat.Unit.ID)": obj.Chat.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		roomID   = obj.RoomID
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   = obj.Chat.Unit.ID
		holder   models.AllowHolder
	)

	if node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(model.ActionTypeRead, roomID, &holder) {
		return node.GetError(), nil
	}
	form, err := r.Services.Repos.Rooms.RoomForm(roomID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	return form, nil
}

func (r *roomResolver) Allows(ctx context.Context, obj *model.Room) (model.AllowsResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Room.Allows", &bson.M{
		"roomID (obj.RoomID)": obj.RoomID,
	})
	defer node.MethodTiming()

	allows, err := r.Services.Repos.Rooms.Allows(obj.RoomID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при попытке получить данные"), nil
	}

	return allows, nil
}

func (r *roomResolver) Messages(ctx context.Context, obj *model.Room, find model.FindMessagesInRoom) (model.MessagesResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Room.Messages", &bson.M{
		"roomID (obj.RoomID)":       obj.RoomID,
		"chatID (obj.Chat.Unit.ID)": obj.Chat.Unit.ID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		roomID   = obj.RoomID
		chatID   = obj.Chat.Unit.ID
		holder   models.AllowHolder
	)

	if node.ValidFindMessagesInRoom(&find) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(model.ActionTypeRead, roomID, &holder) {
		return node.GetError(), nil
	}
	room, err := r.Services.Repos.Messages.MessagesFromRoom(roomID, chatID, &find)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrBadRequest, "произошла ошибка во время обработки данных"), nil
	}

	return room, nil
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// ListenCollection returns generated.ListenCollectionResolver implementation.
func (r *Resolver) ListenCollection() generated.ListenCollectionResolver {
	return &listenCollectionResolver{r}
}

// Me returns generated.MeResolver implementation.
func (r *Resolver) Me() generated.MeResolver { return &meResolver{r} }

// Member returns generated.MemberResolver implementation.
func (r *Resolver) Member() generated.MemberResolver { return &memberResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type chatResolver struct{ *Resolver }
type listenCollectionResolver struct{ *Resolver }
type meResolver struct{ *Resolver }
type memberResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type roomResolver struct{ *Resolver }

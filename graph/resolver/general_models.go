package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *chatResolver) Owner(ctx context.Context, obj *model.Chat) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) Rooms(ctx context.Context, obj *model.Chat) ([]*model.Room, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) CountMembers(ctx context.Context, obj *model.Chat) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) Members(ctx context.Context, obj *model.Chat) ([]*model.ChatMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) Roles(ctx context.Context, obj *model.Chat) ([]*model.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) Invites(ctx context.Context, obj *model.Chat) ([]*model.Invite, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) Banlist(ctx context.Context, obj *model.Chat) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatResolver) MeRestricts(ctx context.Context, obj *model.Chat) ([]*model.MeRestricts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatMemberResolver) Chat(ctx context.Context, obj *model.ChatMember) (*model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *chatMemberResolver) User(ctx context.Context, obj *model.ChatMember) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *inviteResolver) Chat(ctx context.Context, obj *model.Invite) (*model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *meRestrictsResolver) Chat(ctx context.Context, obj *model.MeRestricts) (*model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
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

func (r *permissionHoldersResolver) Roles(ctx context.Context, obj *model.PermissionHolders) ([]*model.Role, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *permissionHoldersResolver) Chars(ctx context.Context, obj *model.PermissionHolders) ([]model.Char, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *permissionHoldersResolver) Members(ctx context.Context, obj *model.PermissionHolders) ([]*model.ChatMember, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roleResolver) Users(ctx context.Context, obj *model.Role) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roomResolver) Chat(ctx context.Context, obj *model.Room) (*model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roomResolver) Restricts(ctx context.Context, obj *model.Room) (*model.Restricts, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *roomResolver) Messages(ctx context.Context, obj *model.Room) ([]*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

// Chat returns generated.ChatResolver implementation.
func (r *Resolver) Chat() generated.ChatResolver { return &chatResolver{r} }

// ChatMember returns generated.ChatMemberResolver implementation.
func (r *Resolver) ChatMember() generated.ChatMemberResolver { return &chatMemberResolver{r} }

// Invite returns generated.InviteResolver implementation.
func (r *Resolver) Invite() generated.InviteResolver { return &inviteResolver{r} }

// MeRestricts returns generated.MeRestrictsResolver implementation.
func (r *Resolver) MeRestricts() generated.MeRestrictsResolver { return &meRestrictsResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// PermissionHolders returns generated.PermissionHoldersResolver implementation.
func (r *Resolver) PermissionHolders() generated.PermissionHoldersResolver {
	return &permissionHoldersResolver{r}
}

// Role returns generated.RoleResolver implementation.
func (r *Resolver) Role() generated.RoleResolver { return &roleResolver{r} }

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type chatResolver struct{ *Resolver }
type chatMemberResolver struct{ *Resolver }
type inviteResolver struct{ *Resolver }
type meRestrictsResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type permissionHoldersResolver struct{ *Resolver }
type roleResolver struct{ *Resolver }
type roomResolver struct{ *Resolver }

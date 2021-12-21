package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) BanMember(ctx context.Context, memberID int) (model.MutationResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) BanUser(ctx context.Context, userID int, chatID int) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.Can.Ban(clientID, chatID) ||
		pl.UserExists(userID) ||
		pl.IsMember(userID, chatID) {
		return pl.Err, nil
	}

	if r.Services.Repos.Chats.AddToBanlist(userID, chatID) != nil ||
		r.Services.Repos.Chats.RemoveUserFromChat(userID, chatID) != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("пользователь успешно забанен"), nil
}

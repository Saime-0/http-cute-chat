package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/api/resp"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/its"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) BanUser(ctx context.Context, userID int, chatID int) (model.MutationResult, error) {
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.UserExists(userID) ||
		pl.ChatExists(chatID) ||
		pl.UserIs(userID, chatID, its.List(its.Member)) {
		return pl.Err, nil
	}
	err := r.Services.Repos.Chats.BanUserInChat(userID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), err
	}
	return resp.Success("пользователь успешно забанен"), nil
}

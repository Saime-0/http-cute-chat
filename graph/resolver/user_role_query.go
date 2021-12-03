package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/its"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) UserRole(ctx context.Context, userID int, chatID int) (model.UserRoleResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.UserIs(chatID, clientID, its.List(its.Member)) ||
		pl.UserExists(userID) ||
		pl.UserIs(chatID, userID, its.List(its.Member)) {
		return pl.Err, nil
	}
	role, err := r.Services.Repos.Chats.UserRole(userID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return role, nil
}

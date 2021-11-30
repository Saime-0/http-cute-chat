package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) GiveRole(ctx context.Context, userID int, chatID int, roleID int) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.Can.GiveRole(clientID, chatID) ||
		pl.RoleExists(chatID, roleID) ||
		pl.UserExists(userID) ||
		pl.IsMember(userID, chatID) {
		return pl.Err, nil
	}

	err := r.Services.Repos.Chats.GiveRole(clientID, chatID, roleID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("роль выдана"), nil
}

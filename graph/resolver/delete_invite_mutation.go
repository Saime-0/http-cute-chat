package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/piping"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) DeleteInvite(ctx context.Context, chatID int, code string) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.Can.CreateInvite(clientID, chatID) ||
		pl.HasInvite(chatID, code) ||
		pl.InviteIsRelevant(code) {
		return pl.Err, nil
	}

	err := r.Services.Repos.Chats.DeleteInviteLinkByCode(code)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("инвайт успешно удален"), nil
}

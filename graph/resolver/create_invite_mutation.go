package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) CreateInvite(ctx context.Context, chatID int, input model.CreateInviteInput) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.Can.CreateInvite(clientID, chatID) ||
		pl.InvitesLimit(chatID) ||
		pl.ValidInviteInput(input) {
		return pl.Err, nil
	}

	_, err := r.Services.Repos.Chats.CreateInviteLink(
		&models.CreateInvite{
			ChatID: chatID,
			Aliens: *input.Aliens,
			Exp:    *input.Duration + time.Now().UTC().Unix(),
		},
	)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("инвайт создан"), nil
}

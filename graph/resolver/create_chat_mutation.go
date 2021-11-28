package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) CreateChat(ctx context.Context, input model.CreateChatInput) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.OwnedLimit(clientID) ||
		pl.ChatCountLimit(clientID) ||
		pl.ValidDomain(input.Domain) ||
		pl.ValidName(input.Name) ||
		pl.DomainIsFree(input.Domain) {
		return pl.Err, nil
	}

	_, err := r.Services.Repos.Chats.CreateChat(clientID,
		&models.CreateChat{
			Domain:  input.Domain,
			Name:    input.Name,
			Private: input.Private,
		},
	)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return resp.Success("чат успешно создан"), nil
}

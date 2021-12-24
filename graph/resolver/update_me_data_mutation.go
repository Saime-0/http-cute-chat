package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/tlog"
)

func (r *mutationResolver) UpdateMeData(ctx context.Context, input model.UpdateMeDataInput) (model.MutationResult, error) {
	tl := tlog.Start("mutationResolver > UpdateMeData [_]")
	defer tl.Fine()

	clientID := ctx.Value(rules.UserIDFromToken).(int)

	node := r.Piper.CreateNode()
	defer node.Kill()

	if input.Name != nil && node.ValidName(*input.Name) ||
		input.Domain != nil && node.ValidDomain(*input.Domain) ||
		input.Password != nil && node.ValidPassword(*input.Password) ||
		input.Email != nil && node.ValidDomain(*input.Email) {
		return node.Err, nil
	}

	err := r.Services.Repos.Users.UpdateMe(clientID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные"), nil
	}

	return resp.Success("данные пользователя обновлены"), nil
}

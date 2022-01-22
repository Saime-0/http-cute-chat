package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *queryResolver) Me(ctx context.Context) (model.MeResult, error) {
	node := r.Piper.CreateNode("queryResolver > Me [_]")
	defer node.Kill()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID
	me, err := r.Services.Repos.Users.Me(clientID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось получить данные"), nil
	}

	return me, nil
}

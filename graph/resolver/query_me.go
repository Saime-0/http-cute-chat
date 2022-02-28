package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *queryResolver) Me(ctx context.Context) (model.MeResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Me", nil)
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID
	me, err := r.Services.Repos.Users.Me(clientID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось получить данные"), nil
	}
	if me == nil {
		return resp.Error(resp.ErrBadRequest, "пользователя не существует"), nil
	}

	return me, nil
}

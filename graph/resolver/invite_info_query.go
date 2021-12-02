package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) InviteInfo(ctx context.Context, code string) (model.InviteInfoResult, error) {
	//clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.InviteIsRelevant(code) {
		return pl.Err, nil
	}

	info, err := r.Services.Repos.Chats.InviteInfo(code)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return info, nil
}

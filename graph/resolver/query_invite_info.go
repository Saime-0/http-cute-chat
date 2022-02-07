package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) InviteInfo(ctx context.Context, code string) (model.InviteInfoResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("InviteInfo", &bson.M{
		"code": code,
	})
	defer node.MethodTiming()

	if node.InviteIsRelevant(code) {
		return node.GetError(), nil
	}

	info, err := r.Services.Repos.Chats.InviteInfo(code)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось получить информацию"), nil
	}
	return info, nil
}

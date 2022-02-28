package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Chats(ctx context.Context, find model.FindChats, params *model.Params) (model.ChatsResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Chats", &bson.M{
		"find":   find,
		"params": params,
	})
	defer node.MethodTiming()

	//node.Debug(cerrors.Wrap(cerrors.New("debuging error testing"), utils.GetCallerPos()))

	if node.ValidParams(&params) ||
		find.ID != nil && node.ValidID(*find.ID) ||
		find.Domain != nil && node.ValidDomain(*find.Domain) ||
		find.NameFragment != nil && node.ValidNameFragment(*find.NameFragment) {
		return node.GetError(), nil
	}

	chats, err := r.Services.Repos.Chats.FindChats(&find, params)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось получить список чатов"), nil
	}

	return chats, nil
}

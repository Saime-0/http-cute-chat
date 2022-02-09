package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Users(ctx context.Context, find model.FindUsers, params *model.Params) (model.UsersResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Users", &bson.M{
		"find":   find,
		"params": params,
	})
	defer node.MethodTiming()

	if node.ValidParams(&params) ||
		find.ID != nil && node.ValidID(*find.ID) ||
		find.Domain != nil && node.ValidDomain(*find.Domain) ||
		find.NameFragment != nil && node.ValidNameFragment(*find.NameFragment) {
		return node.GetError(), nil
	}

	users, err := r.Services.Repos.Users.FindUsers(&find)
	if err != nil {
		node.Healer.Alert(errors.Wrap(err, utils.GetCallerPos()))
		return nil, errors.New("произошла ошибка во время обработки данных")
	}

	return users, nil
}

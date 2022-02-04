package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) Units(ctx context.Context, find model.FindUnits, params *model.Params) (model.UnitsResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Units", &bson.M{
		"find":   find,
		"params": params,
	})
	defer node.MethodTiming()

	var units *model.Units

	if node.ValidParams(&params) ||
		find.ID != nil && node.ValidID(*find.ID) ||
		find.NameFragment != nil && node.ValidNameFragment(*find.NameFragment) ||
		find.Domain != nil && node.ValidNameFragment(*find.Domain) {
		return node.GetError(), nil
	}

	units = r.Services.Repos.Units.FindUnits(&find, params)
	return units, nil
}

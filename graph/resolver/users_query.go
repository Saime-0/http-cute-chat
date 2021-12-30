package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) Users(ctx context.Context, find model.FindUsers, params *model.Params) (model.UsersResult, error) {
	node := r.Piper.CreateNode("queryResolver > Users [_]")
	defer node.Kill()

	if node.ValidParams(&params) ||
		find.ID != nil && node.ValidID(*find.ID) ||
		find.Domain != nil && node.ValidDomain(*find.Domain) ||
		find.NameFragment != nil && node.ValidNameFragment(*find.NameFragment) {
		return node.Err, nil
	}

	users := r.Services.Repos.Users.FindUsers(&find)
	return users, nil
}

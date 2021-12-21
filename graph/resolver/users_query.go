package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Users(ctx context.Context, find model.FindUsers, params *model.Params) (model.UsersResult, error) {
	pl := piping.NewPipeline(r.Services.Repos)
	if pl.ValidParams(params) ||
		find.UserID != nil && pl.ValidID(*find.UserID) ||
		find.UserDomain != nil && pl.ValidDomain(*find.UserDomain) ||
		find.NameFragment != nil && pl.ValidNameFragment(*find.NameFragment) {
		return pl.Err, nil
	}

	users := r.Services.Repos.Users.FindUsers(&find)
	return users, nil
}

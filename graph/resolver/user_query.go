package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) User(ctx context.Context, id *int64, domain *string) (model.UserResult, error) {
	return model.User{Unit: &model.Unit{
		ID:     0,
		Domain: "afh",
		Name:   "hmdg",
		Type:   "fGH",
	}}, nil
}

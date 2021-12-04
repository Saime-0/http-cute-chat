package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) Members(ctx context.Context, chatID int, userID *int) (model.MembersResult, error) {
	panic(fmt.Errorf("not implemented"))
}
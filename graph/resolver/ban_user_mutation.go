package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/pipeline"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/its"
)

func (r *mutationResolver) BanUser(ctx context.Context, userID int64, chatID int64) (model.MutationResult, error) {
	pl := pipeline.NewPipeline(ctx, r.Services.Repos)

	if pl.UserIs(its.Owner, its.Admin) ||
		pl.UserExists() ||
		pl.ChatExists() {

	}

}

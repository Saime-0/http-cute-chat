package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *mutationResolver) JoinByInvite(ctx context.Context, code string) (model.JoinByInviteResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.InviteIsExists(code) || // ?
		pl.InviteIsRelevant(code) ||
		pl.GetChatByInvite(code) ||
		pl.IsNotMember(clientID) {
		return pl.Err, nil
	}
}

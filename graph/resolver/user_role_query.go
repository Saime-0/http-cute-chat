package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) UserRole(ctx context.Context, memberID int) (model.UserRoleResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	var chatID int
	if pl.FindMember(memberID, &chatID) ||
		pl.IsMember(clientID, chatID) {
		return pl.Err, nil
	}
	role := r.Services.Repos.Chats.MemberRole(memberID)

	return role, nil
}

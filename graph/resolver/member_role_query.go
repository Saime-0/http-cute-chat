package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/rules"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) MemberRole(ctx context.Context, memberID int) (model.UserRoleResult, error) {
	node := r.Piper.CreateNode("queryResolver > MemberRole [mid:", memberID, "]")
	defer node.Kill()

	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		chatID   int
	)

	if node.GetChatIDByMember(memberID, &chatID) ||
		node.IsMember(clientID, chatID) {
		return node.Err, nil
	}
	role := r.Services.Repos.Chats.MemberRole(memberID)

	return role, nil
}

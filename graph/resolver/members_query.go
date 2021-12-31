package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *queryResolver) Members(ctx context.Context, find model.FindMembers) (model.MembersResult, error) {
	node := r.Piper.CreateNode("queryResolver > Members [cid:", find.ChatID, "]")
	defer node.Kill()

	var (
		chatID   int
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		members  *model.Members
	)

	if node.ValidID(chatID) ||
		node.IsMember(clientID, chatID) ||
		find.MemberID != nil && node.ValidID(*find.MemberID) ||
		find.RoleID != nil && node.ValidID(*find.RoleID) {
		return node.Err, nil
	}

	members = r.Services.Repos.Chats.FindMembers(&find)
	return members, nil
}

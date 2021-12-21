package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Members(ctx context.Context, find model.FindMembers) (model.MembersResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var (
		chatID  = find.ChatID
		members *model.Members
	)

	if pl.ValidID(chatID) ||
		pl.IsMember(clientID, chatID) ||
		find.MemberID != nil && pl.ValidID(*find.MemberID) ||
		find.RoleID != nil && pl.ValidID(*find.RoleID) {
		return pl.Err, nil
	}

	members = r.Services.Repos.Chats.FindMembers(&find)
	return members, nil
}

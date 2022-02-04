package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *queryResolver) Members(ctx context.Context, find model.FindMembers) (model.MembersResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Members", &bson.M{
		"find": find,
	})
	defer node.MethodTiming()

	var (
		chatID   int
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		members  *model.Members
	)

	if node.ValidID(chatID) ||
		node.IsMember(clientID, chatID) ||
		find.MemberID != nil && node.ValidID(*find.MemberID) ||
		find.RoleID != nil && node.ValidID(*find.RoleID) {
		return node.GetError(), nil
	}

	members = r.Services.Repos.Chats.FindMembers(&find)
	return members, nil
}

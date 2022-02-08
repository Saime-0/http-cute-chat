package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) MemberRole(ctx context.Context, memberID int) (model.UserRoleResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("MemberRole", &bson.M{
		"memberID": memberID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   int
	)

	if node.GetChatIDByMember(memberID, &chatID) ||
		node.IsMember(clientID, chatID) {
		return node.GetError(), nil
	}
	role, _ := r.Dataloader.MemberRole(memberID)

	return role, nil
}

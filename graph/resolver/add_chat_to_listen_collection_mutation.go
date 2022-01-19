package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/models"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) AddChatToListenCollection(ctx context.Context, chatID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > AddChatToListenCollection [cid:", chatID, "]")
	defer node.Kill()

	var (
		clientID    = ctx.Value(rules.UserIDFromToken).(int)
		memberID    int
		userMembers *[]*models.SubUser
	)

	if node.ChatExists(chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) ||
		node.UserHasAccessToChats(clientID, &[]int{chatID}, &userMembers) {
		return node.Err, nil
	}

	//r.Services.Subix.CreateMemberIfNotExists(memberID)
	panic("not implemented")
}

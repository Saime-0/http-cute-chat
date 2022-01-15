package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/rules"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *mutationResolver) AddChatToListenCollection(ctx context.Context, chatID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > AddChatToListenCollection [cid:", chatID, "]")
	defer node.Kill()

	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		memberID int
	)

	if node.ChatExists(chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) {
		return node.Err, nil
	}

	panic("Not implemented")
}

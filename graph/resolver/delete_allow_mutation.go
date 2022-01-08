package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *mutationResolver) DeleteAllow(ctx context.Context, allowID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > DeleteAllow [aid:", allowID, "]")
	defer node.Kill()

	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		chatID   int
	)

	if node.ValidID(allowID) ||
		node.GetChatIDByAllow(allowID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateAllow(clientID, chatID) {
		return node.Err, nil
	}

	eventReadyAllow, err := r.Services.Repos.Rooms.DeleteAllow(allowID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось удалить разрешение"), nil
	}
	r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyAllow,
	)
	return resp.Success("успешно удалено"), nil
}

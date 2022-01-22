package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *mutationResolver) DeleteAllow(ctx context.Context, allowID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > DeleteAllow [aid:", allowID, "]")
	defer node.Kill()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
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
	go r.Services.Subix.NotifyChatMembers(
		chatID,
		eventReadyAllow,
	)
	return resp.Success("успешно удалено"), nil
}

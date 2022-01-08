package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) CreateAllow(ctx context.Context, roomID int, input model.AllowInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > CreateAllow [rid:", roomID, "]")
	defer node.Kill()

	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		chatID   int
	)

	if node.ValidID(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateAllow(clientID, chatID) ||
		node.ValidAllowInput(chatID, &input) ||
		node.AllowNotExists(roomID, &input) {
		return node.Err, nil
	}

	eventReadyAllow, err := r.Services.Repos.Rooms.CreateAllow(roomID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать разрешение"), nil
	}
	r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyAllow,
	)
	return resp.Success("успешно создано"), nil
}

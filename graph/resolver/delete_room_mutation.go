package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) DeleteRoom(ctx context.Context, roomID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > DeleteRoom [rid:", roomID, "]")
	defer node.Kill()
	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		chatID   int
	)

	if node.ValidID(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateRoom(clientID, chatID) {
		return node.Err, nil
	}

	eventReadyRoom, err := r.Services.Repos.Rooms.DeleteRoom(roomID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось удалить комнату"), nil
	}
	r.Services.Subix.NotifyChatMembers(
		[]int{chatID},
		eventReadyRoom,
	)
	return resp.Success("комната удалена"), nil
}

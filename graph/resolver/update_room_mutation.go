package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) UpdateRoom(ctx context.Context, roomID int, input model.UpdateRoomInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > UpdateRoom [rid:", roomID, "]")
	defer node.Kill()

	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		chatID   int
	)

	if node.GetChatIDByRoom(roomID, &chatID) ||
		node.CanUpdateRoom(clientID, chatID) ||
		input.Name != nil && node.ValidName(*input.Name) ||
		input.Note != nil && node.ValidNote(*input.Note) ||
		input.ParentID != nil && node.ValidParentRoomID(roomID, *input.ParentID) {
		return node.Err, nil
	}

	eventReadyRoom, err := r.Services.Repos.Rooms.UpdateRoom(roomID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные комнаты"), nil
	}

	go r.Services.Subix.NotifyChatMembers(
		chatID,
		eventReadyRoom,
	)

	return resp.Success("данные комнаты обновлены"), nil
}

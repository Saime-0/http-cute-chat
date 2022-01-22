package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *mutationResolver) CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > CreateRoom [cid:", input.ChatID, "]")
	defer node.Kill()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(input.ChatID) ||
		node.IsMember(clientID, input.ChatID) ||
		node.CanCreateRoom(clientID, input.ChatID) ||
		node.RoomsLimit(input.ChatID) ||
		input.Parent != nil && node.IsNotChild(*input.Parent) ||
		input.Parent != nil && node.RoomExists(*input.Parent) ||
		input.Allows != nil && node.ValidRoomAllows(input.ChatID, input.Allows) {
		return node.Err, nil
	}

	eventReadyRoom, err := r.Services.Repos.Rooms.CreateRoom(&input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать комнату"), nil
	}
	go r.Services.Subix.NotifyChatMembers(
		input.ChatID,
		eventReadyRoom,
	)
	return resp.Success("комната создана"), nil
}

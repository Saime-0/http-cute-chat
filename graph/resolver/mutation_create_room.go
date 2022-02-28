package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateRoom(ctx context.Context, input model.CreateRoomInput) (model.CreateRoomResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("CreateRoom", &bson.M{
		"input": input,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(input.ChatID) ||
		node.IsMember(clientID, input.ChatID) ||
		node.CanCreateRoom(clientID, input.ChatID) ||
		node.RoomsLimit(input.ChatID) ||
		input.Parent != nil && node.IsNotChild(*input.Parent) ||
		input.Parent != nil && node.RoomExists(*input.Parent) ||
		input.Allows != nil && node.ValidRoomAllows(input.ChatID, input.Allows) {
		return node.GetError(), nil
	}

	eventReadyRoom, err := r.Services.Repos.Rooms.CreateRoom(&input)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось создать комнату"), nil
	}
	go r.Subix.NotifyChatMembers(
		input.ChatID,
		eventReadyRoom,
	)
	return &model.CreatedRoom{RoomID: eventReadyRoom.ID}, nil
}

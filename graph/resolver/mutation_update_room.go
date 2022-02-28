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

func (r *mutationResolver) UpdateRoom(ctx context.Context, roomID int, input model.UpdateRoomInput) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("UpdateRoom", &bson.M{
		"roomID": roomID,
		"input":  input,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   int
	)

	if node.GetChatIDByRoom(roomID, &chatID) ||
		node.CanUpdateRoom(clientID, chatID) ||
		input.Name != nil && node.ValidName(*input.Name) ||
		input.Note != nil && node.ValidNote(*input.Note) ||
		input.ParentID != nil && node.ValidParentRoomID(roomID, *input.ParentID) {
		return node.GetError(), nil
	}

	eventReadyRoom, err := r.Services.Repos.Rooms.UpdateRoom(roomID, &input)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные комнаты"), nil
	}

	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyRoom,
	)

	return resp.Success("данные комнаты обновлены"), nil
}

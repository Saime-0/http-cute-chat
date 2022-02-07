package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateAllows(ctx context.Context, roomID int, input model.AllowsInput) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("ConfirmRegistration", &bson.M{
		"roomID": roomID,
		"input":  input,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   int
	)

	if node.ValidID(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateAllow(clientID, chatID) ||
		node.AllowsNotExists(roomID, &input) ||
		node.ValidRoomAllows(chatID, &input) {
		return node.GetError(), nil
	}

	eventReadyAllow, err := r.Services.Repos.Rooms.CreateAllows(roomID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать разрешение"), nil
	}
	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyAllow,
	)
	return resp.Success("успешно создано"), nil
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *queryResolver) RoomForm(ctx context.Context, roomID int) (model.RoomFormResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		chatID   int
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		holder   models.AllowHolder
	)

	if node.ValidID(roomID) ||
		node.RoomExists(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(model.ActionTypeRead, roomID, &holder) {
		return node.Err, nil
	}

	form := r.Services.Repos.Rooms.RoomForm(roomID)

	return form, nil
}

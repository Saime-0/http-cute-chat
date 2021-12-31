package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *queryResolver) RoomForm(ctx context.Context, roomID int) (model.RoomFormResult, error) {
	node := r.Piper.CreateNode("queryResolver > RoomForm [rid:", roomID, "]")
	defer node.Kill()

	var (
		chatID   int
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		holder   models.AllowHolder
	)

	if node.ValidID(roomID) ||
		node.RoomExists(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(rules.AllowRead, roomID, &holder) {
		return node.Err, nil
	}

	form := r.Services.Repos.Rooms.RoomForm(roomID)

	return form, nil
}

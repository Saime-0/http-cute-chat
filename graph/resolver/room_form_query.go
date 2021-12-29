package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *queryResolver) RoomForm(ctx context.Context, roomID int) (model.RoomFormResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var (
		chatID int
		holder models.AllowHolder
	)
	if pl.ValidID(roomID) ||
		pl.RoomExists(roomID) ||
		pl.GetChatIDByRoom(roomID, &chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.GetAllowHolder(clientID, chatID, &holder) ||
		pl.IsAllowedTo(rules.AllowRead, roomID, &holder) {
		return pl.Err, nil
	}

	form := r.Services.Repos.Rooms.RoomForm(roomID)

	return form, nil
}

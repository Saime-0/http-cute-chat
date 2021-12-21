package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Message(ctx context.Context, msgID int) (model.MessageResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var (
		chatID  int
		holder  models.AllowHolder
		message model.Message
	)
	if pl.GetMessageByID(msgID, &message) ||
		pl.GetChatIDByRoom(message.Room.RoomID, &chatID) ||
		pl.IsMember(clientID, chatID) ||
		(r.Services.Repos.Rooms.AllowsIsSet(message.Room.RoomID)) &&
			(pl.GetAllowHolder(clientID, chatID, &holder) ||
				pl.IsAllowedTo(rules.AllowRead, message.Room.RoomID, &holder)) {
		return pl.Err, nil
	}
	panic("Not Implemented")
	return message, nil
}

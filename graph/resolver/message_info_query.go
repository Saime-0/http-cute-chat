package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) MessageInfo(ctx context.Context, id int, roomID int) (model.MessageInfoResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	var chatID int
	var holder *models.AllowHolder
	if pl.RoomExists(roomID) ||
		pl.GetChatIDByRoom(roomID, &chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.MessageAvailable(id, roomID) ||
		(r.Services.Repos.Rooms.AllowsIsSet(roomID)) &&
			(pl.GetAllowHolder(clientID, chatID, holder) ||
				pl.IsAllowedTo(rules.AllowRead, roomID, holder)) {
		return pl.Err, nil
	}

	message, err := r.Services.Repos.Messages.Message(id)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return message, nil
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *mutationResolver) CreateRoom(ctx context.Context, chatID int, input model.CreateRoomInput) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.Can.CreateRoom(clientID, chatID) ||
		pl.CountRoomLimit(chatID) ||
		input.Parent != nil && pl.RoomExists(*input.Parent) {
		return pl.Err, nil
	}
	_, err := r.Services.Repos.Rooms.CreateRoom(chatID, &models.CreateRoom{
		Name:      input.Name,
		ParentID:  *input.Parent,
		Note:      *input.Note, // ! ptr означает что значение может быть не указано
		MsgFormat: models.FormPattern{},
		Restricts: models.Allows{},
	})
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (r *mutationResolver) UpdateRoomForm(ctx context.Context, roomID int, form *model.UpdateFormInput) (model.MutationResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var chatID int
	if pl.RoomExists(roomID) ||
		pl.GetChatIDByRoom(roomID, &chatID) ||
		pl.IsMember(clientID, chatID) ||
		pl.Can.UpdateRoom(clientID, chatID) ||
		form != nil && pl.ValidForm(form) {
		return pl.Err, nil
	}
	var (
		err      error
		bodyForm *string
	)
	if form != nil {
		byteForm, err := json.Marshal(*form)
		bodyForm = kit.StringPtr(string(byteForm))
		if err != nil {
			return resp.Error(resp.ErrInternalServerError, "не удалось обработать тело запроса"), nil
		}
	}

	err = r.Services.Repos.Rooms.UpdateRoomForm(roomID, bodyForm)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось установить новую форму"), nil
	}
	return resp.Success("форма успешно обновлена"), nil
}

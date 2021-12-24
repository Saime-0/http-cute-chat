package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/tlog"
)

func (r *mutationResolver) UpdateRoom(ctx context.Context, roomID int, input model.UpdateRoomInput) (model.MutationResult, error) {
	tl := tlog.Start("mutationResolver > UpdateRoom [rid:", roomID, "]")
	defer tl.Fine()

	clientID := ctx.Value(rules.UserIDFromToken).(int)

	node := r.Piper.CreateNode()
	defer node.Kill()

	var (
		chatID int
	)

	if node.GetChatIDByRoom(roomID, &chatID) ||
		node.CanUpdateRoom(clientID, chatID) ||
		input.Name != nil && node.ValidName(*input.Name) ||
		input.Note != nil && node.ValidNote(*input.Note) ||
		input.ParentID != nil && node.ValidID(*input.ParentID) {
		return node.Err, nil
	}

	err := r.Services.Repos.Rooms.UpdateRoom(roomID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные комнаты"), nil
	}

	return resp.Success("данные комнаты обновлены"), nil
}

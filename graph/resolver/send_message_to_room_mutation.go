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

func (r *mutationResolver) SendMessageToRoom(ctx context.Context, roomID int, input model.CreateMessageInput) (model.SendMessageToRoomResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var (
		chatID   int
		memberID int
		holder   models.AllowHolder
		//handledUserChoice string
	)
	if pl.RoomExists(roomID) ||
		pl.GetChatIDByRoom(roomID, &chatID) ||
		pl.GetMemberBy(clientID, chatID, &memberID) ||
		pl.IsNotMuted(memberID) ||
		pl.GetAllowHolder(clientID, chatID, &holder) ||
		pl.IsAllowedTo(rules.AllowWrite, roomID, &holder) ||
		r.Services.Repos.Rooms.FormIsSet(roomID) && pl.HandleChoice(input.Body, roomID, &input.Body) {
		return pl.Err, nil
	}
	// todo if message is anonimus or room
	// todo все allow.group.user заменить на member
	message := &models.CreateMessage{
		ReplyTo: input.ReplyTo,
		Author:  &clientID,
		RoomID:  roomID,
		Body:    input.Body,
		Type:    model.MessageTypeUser,
	}
	err := r.Services.Repos.Messages.CreateMessageInRoom(message)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать сообщение"), nil
	}
	return resp.Success("сообщение успешно создано"), nil
}

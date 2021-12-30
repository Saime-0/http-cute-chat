package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) SendMessageToRoom(ctx context.Context, roomID int, input model.CreateMessageInput) (model.SendMessageToRoomResult, error) {
	node := r.Piper.CreateNode("mutationResolver > SendMessageToRoom [rid:", roomID, "]")
	defer node.Kill()
	var (
		chatID   int
		memberID int
		holder   models.AllowHolder
		//handledUserChoice string
		clientID = ctx.Value(rules.UserIDFromToken).(int)
	)

	if node.RoomExists(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) ||
		node.IsNotMuted(memberID) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(rules.AllowWrite, roomID, &holder) ||
		r.Services.Repos.Rooms.FormIsSet(roomID) && node.HandleChoice(input.Body, roomID, &input.Body) {
		return node.Err, nil
	}
	println(clientID, "clientID\n", memberID, "memberID\n", chatID, "chatID\n")
	// todo if message is anonimus or room
	// todo все allow.group.user заменить на member
	message := &models.CreateMessage{
		ReplyTo: input.ReplyTo,
		Author:  &memberID,
		RoomID:  roomID,
		Body:    input.Body,
		Type:    model.MessageTypeUser,
	}
	msgID, err := r.Services.Repos.Messages.CreateMessageInRoom(message)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать сообщение"), nil
	}

	var _replyTo *model.Message
	if message.ReplyTo != nil {
		_replyTo = new(model.Message)
	}
	r.Services.Events.NewMessage(roomID, &model.Message{
		ID:      msgID,
		ReplyTo: _replyTo,
		Author:  &model.Member{ID: memberID},
		Type:    message.Type,
	})

	return resp.Success("сообщение успешно создано"), nil
}

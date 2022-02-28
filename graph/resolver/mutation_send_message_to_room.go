package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) SendMessageToRoom(ctx context.Context, roomID int, input model.CreateMessageInput) (model.SendMessageToRoomResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("SendMessageToRoom", &bson.M{
		"roomID": roomID,
		"input":  input,
	})
	defer node.MethodTiming()

	var (
		chatID   int
		memberID int
		holder   models.AllowHolder
		//handledUserChoice string
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if node.RoomExists(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) ||
		node.IsNotMuted(memberID) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(model.ActionTypeWrite, roomID, &holder) ||
		r.Services.Repos.Rooms.FormIsSet(roomID) && node.HandleChoice(input.Body, roomID, &input.Body) {
		return node.GetError(), nil
	}
	// todo if message is anonimus or room
	// todo все allow.group.user заменить на member
	message := &models.CreateMessage{
		ReplyTo: input.ReplyTo,
		UserID:  &clientID,
		RoomID:  roomID,
		Body:    input.Body,
		Type:    model.MessageTypeUser,
	}

	eventReadyMessage, err := func(n piper.Node) (*model.NewMessage, error) {
		n.SwitchMethod("CreateMessageInRoom", &bson.M{
			"ReplyTo": message.ReplyTo,
			"UserID":  message.UserID,
			"RoomID":  message.RoomID,
			"Body":    message.Body,
			"Type":    message.Type,
		})
		defer n.MethodTiming()

		return r.Services.Repos.Messages.CreateMessageInRoom(message)
	}(node)

	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось создать сообщение"), nil
	}

	//r.Services.Events.NewMessage(roomID, &model.Message{ID:      msgID, ReplyTo: _replyTo, UserID:  &model.Member{ID: memberID}, Type:    message.Type, Body:    input.Body})
	go func() {
		err := r.Subix.NotifyRoomReaders(
			roomID,
			eventReadyMessage,
		)
		if err != nil {
			node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		}
	}()

	return resp.Success("сообщение успешно создано"), nil
}

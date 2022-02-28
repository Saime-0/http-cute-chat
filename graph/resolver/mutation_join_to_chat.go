package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) JoinToChat(ctx context.Context, chatID int) (model.JoinToChatResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("JoinToChat", &bson.M{
		"chatID": chatID,
	})
	defer node.MethodTiming()

	var clientID = utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(chatID) ||
		node.IsNotMember(clientID, chatID) ||
		node.IsNotBanned(clientID, chatID) ||
		node.ChatIsNotPrivate(chatID) ||
		node.MembersLimit(chatID) ||
		node.ChatsLimit(clientID) {
		return node.GetError(), nil
	}
	eventReadyMember, err := r.Services.Repos.Chats.AddUserToChat(clientID, chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось присоединиться"), nil
	}

	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyMember,
	)

	return &model.JoinedToChat{ChadID: chatID}, nil
}

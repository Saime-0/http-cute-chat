package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) LeaveFromChat(ctx context.Context, chatID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("LeaveFromChat", &bson.M{
		"chatID": chatID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		memberID int
	)

	if node.ChatExists(chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) ||
		node.CanLeaveFromChat(clientID, chatID) {
		return node.GetError(), nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.RemoveUserFromChat(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	r.Subix.DeleteMember(memberID)
	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyMember,
	)

	return resp.Success("успешно покинул чат"), nil
}

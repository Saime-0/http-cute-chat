package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) TakeChar(ctx context.Context, memberID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("TakeChar", &bson.M{
		"memberID": memberID,
	})
	defer node.MethodTiming()

	var (
		clientID       = utils.GetAuthDataFromCtx(ctx).UserID
		chatID         int
		clientMemberID int
	)

	if node.ValidID(memberID) ||
		node.GetChatIDByMember(memberID, &chatID) ||
		node.GetMemberBy(clientID, chatID, &clientMemberID) ||
		node.CanTakeChar(clientMemberID, memberID) {
		return node.GetError(), nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.TakeChar(memberID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось забрать чар"), nil
	} else if eventReadyMember.ID == 0 {
		return resp.Error(resp.ErrInternalServerError, "у пользователя нет чара"), nil
	}
	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyMember,
	)
	return resp.Success("успешно"), nil
}

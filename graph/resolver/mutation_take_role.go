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

func (r *mutationResolver) TakeRole(ctx context.Context, memberID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("TakeRole", &bson.M{
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
		node.CanTakeRole(clientMemberID, memberID) {
		return node.GetError(), nil
	}

	eventReadyMember, err := r.Services.Repos.Chats.TakeRole(memberID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	} else if eventReadyMember.ID == 0 {
		return resp.Error(resp.ErrInternalServerError, "у пользователя нет роли"), nil
	}
	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyMember,
	)
	return resp.Success("успешно"), nil
}

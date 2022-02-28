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

func (r *mutationResolver) DeleteAllow(ctx context.Context, allowID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("CreateRoom", &bson.M{
		"allowID": allowID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   int
	)

	if node.ValidID(allowID) ||
		node.GetChatIDByAllow(allowID, &chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateAllow(clientID, chatID) {
		return node.GetError(), nil
	}

	eventReadyAllow, err := r.Services.Repos.Rooms.DeleteAllow(allowID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось удалить разрешение"), nil
	}
	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyAllow,
	)
	return resp.Success("успешно удалено"), nil
}

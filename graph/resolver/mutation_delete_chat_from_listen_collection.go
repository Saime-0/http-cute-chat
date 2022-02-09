package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/pkg/errors"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) DeleteChatFromListenCollection(ctx context.Context, sessionKey string, chatID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("DeleteChatFromListenCollection", &bson.M{
		"sessionKey": sessionKey,
		"chatID":     chatID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		memberID int
	)

	if node.ValidSessionKey(sessionKey) ||
		node.ValidID(chatID) ||
		node.ChatExists(chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) {
		return node.GetError(), nil
	}

	err := r.Subix.DeleteChatFromListenCollection(sessionKey, memberID)
	if err != nil {
		node.Healer.Alert(errors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrBadRequest, "не удалось прекратить прослушивать чат"), nil
	}

	return resp.Success("успешно"), nil
}

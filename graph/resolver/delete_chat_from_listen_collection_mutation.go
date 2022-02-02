package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) DeleteChatFromListenCollection(ctx context.Context, sessionKey string, chatID int) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > DeleteChatFromListenCollection [cid:", chatID, "]")
	defer node.Kill()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		memberID int
	)

	if node.ValidSessionKey(sessionKey) ||
		node.ValidID(chatID) ||
		node.ChatExists(chatID) ||
		node.GetMemberBy(clientID, chatID, &memberID) {
		return node.Err, nil
	}

	err := r.Subix.DeleteChatFromListenCollection(sessionKey, memberID)
	if err != nil {
		println("DeleteChatFromListenCollection:", err) // debug
		return resp.Error(resp.ErrBadRequest, "не удалось прекратить прослушивать чат"), nil
	}

	return resp.Success("успешно"), nil
}

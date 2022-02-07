package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) AddChatToListenCollection(ctx context.Context, sessionKey string, chatID int) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("AddChatToListenCollection", &bson.M{
		"sessionKey": sessionKey,
		"chatID":     chatID,
	})
	defer node.MethodTiming()

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		subuser  = &models.SubUser{
			MemberID: kit.IntPtr(0),
			ChatID:   &chatID,
		}
	)

	if node.ValidSessionKey(sessionKey) ||
		node.ChatExists(chatID) ||
		node.GetMemberBy(clientID, chatID, subuser.MemberID) {
		return node.GetError(), nil
	}

	err := r.Subix.AddListenChat(sessionKey, subuser)
	if err != nil {
		println("AddChatToListenCollection:", err) // debug
		return resp.Error(resp.ErrBadRequest, "не удалось начать прослушивать чат"), nil
	}

	return resp.Success("успешно"), nil
}

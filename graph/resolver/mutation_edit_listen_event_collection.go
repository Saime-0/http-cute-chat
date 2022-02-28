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

func (r *mutationResolver) EditListenEventCollection(ctx context.Context, sessionKey string, action model.EventSubjectAction, targetChats []int, listenEvents []model.EventType) (model.EditListenEventCollectionResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("EditListenEventCollection", &bson.M{
		"sessionKey":   sessionKey,
		"action":       action,
		"targetChats":  targetChats,
		"listenEvents": listenEvents,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID
	var (
		userMembers *[]*models.SubUser
	)
	if len(listenEvents) > len(model.AllEventType) {
		return resp.Error(resp.ErrBadRequest, "недопустимая длина списка событий"), nil
	}
	targetChats = kit.GetUniqueInts(targetChats) // избавляемся от повторяющихся значений

	if node.ValidSessionKey(sessionKey) ||
		node.UserHasAccessToChats(clientID, &targetChats, &userMembers) {
		return node.GetError(), nil
	}

	err := r.Subix.ModifyCollection(sessionKey, *userMembers, action, listenEvents)
	if err != nil {
		//node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrBadRequest, "не удалось обновить коллекцию"), nil
	}

	return &model.ListenCollection{SessionKey: sessionKey}, nil
}

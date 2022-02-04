package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (r *subscriptionResolver) NewEvent(ctx context.Context, sessionKey string, listenChatCollection []int) (<-chan *model.SubscriptionBody, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("NewEvent", &bson.M{
		"sessionKey":           sessionKey,
		"listenChatCollection": listenChatCollection,
	})
	defer node.MethodTiming()

	var (
		authData    = utils.GetAuthDataFromCtx(ctx)
		userMembers *[]*models.SubUser
	)

	if authData == nil { // тк @isAuth  вебсокетинге не отрабатывает
		return nil, errors.New("не аутентифицирован")
	}

	listenChatCollection = kit.GetUniqueInts(listenChatCollection) // избавляемся от повторяющихся значений
	if node.ValidSessionKey(sessionKey) ||
		node.UserHasAccessToChats(authData.UserID, &listenChatCollection, &userMembers) {
		return nil, errors.New(node.GetError().Error)
	}

	client, err := r.Subix.Sub(
		authData.UserID,
		sessionKey,
		authData.ExpiresAt,
		*userMembers,
	)
	if err != nil {
		return nil, err
	}

	println("New client", sessionKey) // debug
	go func() {
		<-ctx.Done()
		println("client", sessionKey, "is down.") // debug
		r.Subix.Unsub(sessionKey)
	}()

	return client.Ch, nil
}

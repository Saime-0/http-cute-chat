package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *subscriptionResolver) Subscribe(ctx context.Context, sessionKey string) (<-chan *model.SubscriptionBody, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Subscribe", &bson.M{
		"sessionKey": sessionKey,
	})
	defer node.MethodTiming()

	var (
		authData = utils.GetAuthDataFromCtx(ctx)
	)

	if authData == nil { // тк @isAuth  вебсокетинге не отрабатывает
		node.Debug("не аутентифицирован")
		return nil, cerrors.New("не аутентифицирован")
	}

	if node.ValidSessionKey(sessionKey) {
		return nil, cerrors.New(node.GetError().Error)
	}

	client, err := r.Subix.Sub(
		authData.UserID,
		sessionKey,
		authData.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	// New client
	go func() {
		<-ctx.Done()
		// client is down
		err = r.Subix.Unsub(sessionKey)
		if err != nil {
			node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		}
	}()

	return client.Ch, nil
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *subscriptionResolver) NewEvent(ctx context.Context) (<-chan *model.SubscriptionBody, error) {
	node := r.Piper.CreateNode("subscriptionResolver > NewEvent [_]")
	defer node.Kill()

	if ctx.Value(rules.UserIDFromToken).(int) == 0 {
		return nil, errors.New("не аутентифицирован")
	}
	clientID := ctx.Value(rules.UserIDFromToken).(int)

	client := r.Services.Subix.Register(clientID)
	println("New client chan", client) // debug
	go func() {
		<-ctx.Done()
		println("client chan", client, "is down.") // debug
		r.Services.Subix.Unsubscribe(&client)
	}()

	return client.Ch, nil
}

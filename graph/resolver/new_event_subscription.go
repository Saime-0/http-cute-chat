package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/saime-0/http-cute-chat/internal/models"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *subscriptionResolver) NewEvent(ctx context.Context, listenChatCollection []int) (<-chan *model.SubscriptionBody, error) {
	node := r.Piper.CreateNode("subscriptionResolver > NewEvent [_]")
	defer node.Kill()
	var (
		clientID    = ctx.Value(rules.UserIDFromToken).(int)
		userMembers *[]*models.SubUser
	)
	if clientID == 0 { // тк @isAuth  вебсокетинге не отрабатывает
		return nil, errors.New("не аутентифицирован")

	}
	if node.UserHasAccessToChats(clientID, &listenChatCollection, &userMembers) {
		return nil, errors.New(node.Err.Error)
	}

	client := r.Services.Subix.Sub(clientID, *userMembers)
	println("New client chan", &client) // debug
	go func() {
		<-ctx.Done()
		println("client chan", &client, "is down.") // debug
		r.Services.Subix.Unsub(&client)
	}()

	return client.Ch, nil
}

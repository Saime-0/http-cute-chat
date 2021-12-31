package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *subscriptionResolver) NewRoomMessage(ctx context.Context, roomID int) (<-chan *model.Message, error) {
	node := r.Piper.CreateNode("subscriptionResolver > NewRoomMessage [rid:", roomID, "]")
	defer node.Kill()

	var (
		clientID = ctx.Value(rules.UserIDFromToken).(int)
		chatID   int
		holder   models.AllowHolder
	)
	fmt.Printf("userID: %d\nchatID: %d\n", clientID, chatID) // debug

	if node.ValidID(roomID) ||
		node.GetChatIDByRoom(roomID, &chatID) ||
		node.GetAllowHolder(clientID, chatID, &holder) ||
		node.IsAllowedTo(rules.AllowRead, roomID, &holder) {
		//println(node.Err.Error)
		return nil, errors.New(node.Err.Error)
	}

	sub := r.Services.Events.OnNewMessage.Register(roomID, make(chan *model.Message))
	println("New sub chan", sub) // debug
	//listener := r.Services.Events.SubscribeOnNewMessage(clientID, roomID)
	go func() {
		<-ctx.Done()
		println("sub chan", sub, "is down.") // debug
		r.Services.Events.Unsubscribe(&sub)
	}()

	return sub.Ch.(chan *model.Message), nil
}

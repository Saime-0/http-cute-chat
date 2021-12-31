package mutexmaps

import "github.com/saime-0/http-cute-chat/graph/model"

// deprecated
func (e *EventHandler) SubscribeOnNewMessage(roomID int) **Subscriber {
	sub := e.OnNewMessage.Register(roomID, make(chan *model.Message))

	return &sub
}

package mutexmaps

import (
	"github.com/saime-0/http-cute-chat/graph/model"
)

type Event map[int]*subGroup

type EventHandler struct {
	OnNewMessage *Event // todo > type event map (with method add listener)
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		OnNewMessage: &Event{},
	}
}

func (e *EventHandler) NewMessage(roomID int, message *model.Message) {
	group := e.OnNewMessage.getGroup(roomID)
	group.acrossGroup(func(sub *Subscriber) {
		select {
		case sub.Ch.(chan *model.Message) <- message:
			println("Message write to sub chan", sub) // debug
		default:
			println("sub chan", sub, "is close.. skip") // debug
		}
	})
}

func (g subGroup) acrossGroup(f func(sub *Subscriber)) {
	next := g.Root.next
	for next != nil {
		f(next)
		next = next.next
	}
}

package mutexmaps

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"sync"
)

type Listener struct {
	Ch     chan *model.Message
	userID int
	roomID int
}
type RoomListeners struct {
	Listeners map[int]*Listener
	mu        *sync.Mutex
}

func (e *Events) Unsubscribe(listener *Listener) {
	e.NewMessageEvent[listener.roomID].mu.Lock()
	delete(e.NewMessageEvent[listener.roomID].Listeners, listener.userID)
	e.NewMessageEvent[listener.roomID].mu.Unlock()
}

//type MutexRoom struct {
//	Mu *sync.Mutex
//	Ch chan *model.Message
//}

type Events struct {
	NewMessageEvent map[int]*RoomListeners
}

func NewEvents() *Events {
	return &Events{
		NewMessageEvent: map[int]*RoomListeners{},
	}
}

func (e *Events) SubscribeOnNewMessage(userID, roomID int) *Listener {
	room, ok := e.NewMessageEvent[roomID]
	if !ok {
		room = &RoomListeners{
			Listeners: map[int]*Listener{},
			mu:        new(sync.Mutex),
		}
		e.NewMessageEvent[roomID] = room
		fmt.Println("Listeners created by subscriber", room)
	}
	l, ok := room.Listeners[userID]
	if !ok {
		l = &Listener{
			Ch:     make(chan *model.Message),
			userID: userID,
			roomID: roomID,
		}
		room.Listeners[userID] = l
	}

	fmt.Printf("Listeners: %#v\nListener: %#v\n", room.Listeners, l) // debug
	return l
}

func (e *Events) NewMessage(roomID int, message *model.Message) {
	room, ok := e.NewMessageEvent[roomID]
	if !ok {
		room = &RoomListeners{
			Listeners: map[int]*Listener{},
			mu:        new(sync.Mutex),
		}
		e.NewMessageEvent[roomID] = room
		fmt.Println("Listeners created by message", room)
		return
	}
	fmt.Printf("RoomListeners: %#v\n", room.Listeners) // debug
	for _, listener := range room.Listeners {
		print("Write to ...")
		if listener == nil {
			println("skip.")
		}
		fmt.Println(listener)
		(*listener).Ch <- message
	}
}

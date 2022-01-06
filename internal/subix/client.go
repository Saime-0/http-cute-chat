package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"sync"
)

type version = int

type User struct {
	ID   int
	Root *Client
	mu   *sync.Mutex
}

type Client struct {
	Ch chan *model.SubscriptionBody
	//rev  map[model.EventType]map[int]version
	user *User
	next *Client
	prev *Client
}

func (s *Subscription) Register(userID int) *Client {

	user, ok := s.users[userID]
	if !ok {
		user = &User{
			ID:   userID,
			mu:   new(sync.Mutex),
			Root: &Client{},
		}
		s.users[userID] = user
	}

	user.mu.Lock()
	defer user.mu.Unlock()
	client := &Client{
		Ch: make(chan *model.SubscriptionBody),
		//rev:  map[model.EventType]map[int]version{},
		user: user,
		next: user.Root.next,
		prev: user.Root,
	}

	if user.Root.next != nil {
		user.Root.next.prev = client
	}
	user.Root.next = client

	return client
}

func (s *Subscription) Unsubscribe(client **Client) {
	if client == nil || *client == nil {
		return
	}
	(*client).user.mu.Lock()
	defer (*client).user.mu.Unlock()

	(*client).prev.next = (*client).next
	if (*client).next != nil {
		(*client).next.prev = (*client).prev
	}
	*client = nil // gc moment
}

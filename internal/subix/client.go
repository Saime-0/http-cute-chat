package subix

import (
	"github.com/saime-0/http-cute-chat/graph/model"
)

type Users map[ID]*User

type User struct {
	ID         int
	rootClient *Client
}

type Client struct {
	Ch     chan *model.SubscriptionBody
	UserID int
	next   **Client
}

func (s *Subix) CreateUserIfNotExists(userID int) *User {
	user, ok := s.users[userID]
	if !ok {
		user = &User{
			ID:         userID,
			rootClient: &Client{},
		}
		s.users[userID] = user
		println("Создан user id", userID) // debug
	}
	return user
}

func (s *Subix) deleteUser(userID int) {
	delete(s.users, userID)
	println("удален пользователь с  id =", userID) // debug
	s.deleteMemberByUserID(userID)
}

func (s *Subix) deleteClient(client **Client) {
	if (*client).next == nil {
		s.deleteUser((*client).UserID)
		return
	}
	println("удален клиент", *client) // debug
	*client = *(*client).next
}

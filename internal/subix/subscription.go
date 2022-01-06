package subix

import (
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Subscription struct {
	users map[int]*User
	repo  *repository.Repositories
}

func NewSubscription(repo *repository.Repositories) *Subscription {
	return &Subscription{
		users: map[int]*User{},
		repo:  repo,
	}
}

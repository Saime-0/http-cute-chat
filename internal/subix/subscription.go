package subix

import (
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
)

type Subscription struct {
	users map[int]*User
	repo  *repository.Repositories
	sched *scheduler.Scheduler
	Store *Store
}

func NewSubscription(repo *repository.Repositories, sched *scheduler.Scheduler) *Subscription {
	return &Subscription{
		users: map[int]*User{},
		repo:  repo,
		sched: sched,
		Store: newStore(),
	}
}

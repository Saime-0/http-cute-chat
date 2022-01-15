package subix

import (
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
)

type Subix struct {
	chats   Chats
	members Members
	users   Users

	repo  *repository.Repositories
	sched *scheduler.Scheduler
	Store *Store
}

func NewSubix(repo *repository.Repositories, sched *scheduler.Scheduler) *Subix {
	return &Subix{
		users:   Users{},
		chats:   Chats{},
		members: Members{},
		repo:    repo,
		sched:   sched,
		Store:   newStore(),
	}
}

package subix

import (
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
)

type Subix struct {
	users map[int]*User
	repo  *repository.Repositories
	sched *scheduler.Scheduler
	Store *Store
}

func NewSubix(repo *repository.Repositories, sched *scheduler.Scheduler) *Subix {
	return &Subix{
		users: map[int]*User{},
		repo:  repo,
		sched: sched,
		Store: newStore(),
	}
}

package service

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
	"github.com/saime-0/http-cute-chat/internal/subix"
)

type Services struct {
	Repos     *repository.Repositories
	Subix     *subix.Subix
	Scheduler *scheduler.Scheduler
}

func NewServices(db *sql.DB) *Services {
	s := &Services{
		Repos:     repository.NewRepositories(db),
		Scheduler: scheduler.NewScheduler(),
	}
	s.Subix = subix.NewSubix(s.Repos, s.Scheduler)

	err := s.prepareScheduleInvites()
	if err != nil {
		panic(err)
	}
	return s
}

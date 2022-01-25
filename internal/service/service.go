package service

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/cache"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/email"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
	"github.com/saime-0/http-cute-chat/internal/subix"
)

type Services struct {
	Repos     *repository.Repositories
	Subix     *subix.Subix
	Scheduler *scheduler.Scheduler
	SMTP      *email.SMTPSender
	Cache     *cache.Cache
}

func NewServices(db *sql.DB, cfg *config.Config) *Services {
	s := &Services{
		Repos: repository.NewRepositories(db),
		//Subix:
		Scheduler: scheduler.NewScheduler(),
		SMTP: email.NewSMTPSender(
			cfg.SMTP.Author,
			cfg.SMTP.From,
			cfg.SMTP.Passwd,
			cfg.SMTP.Host,
			cfg.SMTP.Port,
		),
		Cache: cache.NewCache(),
	}
	s.Subix = subix.NewSubix(s.Repos, s.Scheduler)

	err := s.regularSchedule(rules.DurationOfScheduleInterval)
	if err != nil {
		panic(err)
	}

	return s
}

package service

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/email"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/scheduler"
	"github.com/saime-0/http-cute-chat/internal/subix"
)

type Services struct {
	Repos     *repository.Repositories
	Subix     *subix.Subix
	Scheduler *scheduler.Scheduler
	SMTP      *email.SMTPSender
}

func NewServices(db *sql.DB, cfg *config.Config) *Services {
	s := &Services{
		Repos:     repository.NewRepositories(db),
		Scheduler: scheduler.NewScheduler(),
		SMTP: email.NewSMTPSender(
			cfg.SMTP.Author,
			cfg.SMTP.From,
			cfg.SMTP.Passwd,
			cfg.SMTP.Host,
			cfg.SMTP.Port,
		),
	}
	s.Subix = subix.NewSubix(s.Repos, s.Scheduler)

	err := s.prepareScheduleInvites()
	if err != nil {
		panic(err)
	}
	return s
}

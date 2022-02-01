package service

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/cache"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/email"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/subix"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
)

type Services struct {
	Repos     *repository.Repositories
	Subix     *subix.Subix
	Scheduler *scheduler.Scheduler
	SMTP      *email.SMTPSender
	Cache     *cache.Cache
	Logger    *clog.Clog
}

func NewServices(db *sql.DB, cfg *config.Config, logger *clog.Clog) (*Services, error) {
	var err error
	newRepos := repository.NewRepositories(db)
	newSched := scheduler.NewScheduler()
	newSMTPSender, err := email.NewSMTPSender(
		cfg.SMTP.Author,
		cfg.SMTP.From,
		cfg.SMTP.Passwd,
		cfg.SMTP.Host,
		cfg.SMTP.Port,
	)
	if err != nil {
		return nil, err
	}

	newCache := cache.NewCache()
	newSubix := subix.NewSubix(newRepos, newSched)

	s := &Services{
		Repos:     newRepos,
		Subix:     newSubix,
		Scheduler: newSched,
		SMTP:      newSMTPSender,
		Cache:     newCache,
		Logger:    logger,
	}
	err = s.regularSchedule(rules.DurationOfScheduleInterval)
	if err != nil {
		return nil, err
	}

	return s, nil
}

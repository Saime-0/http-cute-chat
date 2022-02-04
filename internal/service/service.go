package service

import (
	"github.com/saime-0/http-cute-chat/internal/cache"
	"github.com/saime-0/http-cute-chat/internal/email"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
)

type Services struct {
	Repos     *repository.Repositories
	Scheduler *scheduler.Scheduler
	SMTP      *email.SMTPSender
	Cache     *cache.Cache
}

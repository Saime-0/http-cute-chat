package service

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/mutexmaps"
	"github.com/saime-0/http-cute-chat/internal/subix"

	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Services struct {
	Repos  *repository.Repositories
	Events *mutexmaps.EventHandler
	Subix  *subix.Subscription
}

func NewServices(db *sql.DB) *Services {
	service := &Services{
		Repos:  repository.NewRepositories(db),
		Events: mutexmaps.NewEventHandler(),
	}
	service.Subix = subix.NewSubscription(service.Repos)
	return service
}

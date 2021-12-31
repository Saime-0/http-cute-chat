package service

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/mutexmaps"

	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Services struct {
	Repos  *repository.Repositories
	Events *mutexmaps.EventHandler
}

func NewServices(db *sql.DB) *Services {
	return &Services{
		Repos:  repository.NewRepositories(db),
		Events: mutexmaps.NewEventHandler(),
	}
}

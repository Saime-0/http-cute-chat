package service

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Services struct {
	Repos *repository.Repositories
}

func NewServices(db *sql.DB) *Services {
	return &Services{
		Repos: repository.NewRepositories(db),
	}
}

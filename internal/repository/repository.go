package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type Users interface {
	Create(u models.User) error
	GetByDomain(u *models.UserDomain) (user models.User, err error)
	Update(inp models.UpdateUserInput) error
	Delete(id int) error
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}

package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(u models.User) (id int, err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name) 
			VALUES ($1, $2) 
			RETURNING id
			) 
		INSERT INTO users (id, app_settings) 
		SELECT u.id, $3 FROM u 
		RETURNING id`,
		u.Domain,
		u.Name,
		u.AppSettings,
	).Scan(&id)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) GetByDomain(domain string) (user models.User, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name,users.app_settings 
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1`,
		domain,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
		&user.AppSettings,
	)
	if err != nil {
		return // user, err
	}

	return
}

func (r *UsersRepo) Update(inp models.UpdateUserInput) error {
	return nil
}
func (r *UsersRepo) Delete(id int) error {
	return nil
}

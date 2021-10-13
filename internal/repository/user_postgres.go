package repository

import (
	"database/sql"
	"encoding/json"
	"log"

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

func (r *UsersRepo) Create(u models.User) error {
	user_json, _ := json.MarshalIndent(u, "", "  ")
	log.Println(string(user_json))
	err := r.db.QueryRow(
		"WITH u AS (INSERT INTO units (domain, name) VALUES ($1, $2) RETURNING id) INSERT INTO users (id, app_settings) SELECT u.id, $3 FROM u",
		u.Domain,
		u.Name,
		u.AppSettings,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) GetByDomain(u *models.UserDomain) (user models.User, err error) {
	// ? select * from units inner join users on units.id = users.id where units.domain = '$1';
	err = r.db.QueryRow(
		"select units.id,units.domain,units.name,users.app_settings from units inner join users on units.id = users.id where units.domain = '$1'",
		u.Domain,
	).Scan(&user)
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

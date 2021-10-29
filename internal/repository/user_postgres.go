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

func (r *UsersRepo) CreateUser(user_model *models.CreateUser) (id int, err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name) 
			VALUES ($1, $2) 
			RETURNING id
			) 
		INSERT INTO users (id, app_settings) 
		SELECT u.id, 'default' FROM u 
		RETURNING id`,
		user_model.Domain,
		user_model.Name,
	).Scan(&id)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) GetUserData(user_id int) (user models.UserData, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id = $1`,
		user_id,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) UserExistsByInput(input_model *models.UserInput) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM units INNER JOIN users 
			ON units.id = users.id 
			WHERE units.domain = $1 AND units.name = $2
			)`,
		input_model.Domain,
		input_model.Name,
	).Scan(&exists)

	return

}

func (r *UsersRepo) GetUserIdByInput(input_model *models.UserInput) (id int, err error) {
	err = r.db.QueryRow(
		`SELECT units.id
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1 AND units.name = $2`,
		input_model.Domain,
		input_model.Name,
	).Scan(&id)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) GetUserByDomain(domain string) (user models.UserInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1`,
		domain,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetUserByID(id int) (user models.UserInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Domain,
		&user.Name,
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetUsersByNameFragment(fragment string, offset int) (users models.ListUserInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id, units.domain,units.name
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.name ILIKE $1
		LIMIT 20
		OFFSET $2`,
		"%"+fragment+"%",
		offset,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.UserInfo{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Name); err != nil {
			return
		}
		users.Users = append(users.Users, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *UsersRepo) GetUserSettings(user_id int) (settings models.UserSettings, err error) {
	err = r.db.QueryRow(
		`SELECT users.app_settings
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id = $1`,
		user_id,
	).Scan(
		&settings.AppSettings,
	)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) UpdateUserData(user_id int, user_model *models.UpdateUserData) error {
	if user_model.Domain != "" {
		err := r.db.QueryRow(
			`UPDATE units
			SET domain = $2
			WHERE id = $1`,
			user_id,
			user_model.Domain,
		).Err()
		if err != nil {
			return err
		}
	}
	if user_model.Name != "" {
		err := r.db.QueryRow(
			`UPDATE units
			SET name = $2
			WHERE id = $1`,
			user_id,
			user_model.Name,
		).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UsersRepo) UpdateUserSettings(user_id int, settings_model *models.UpdateUserSettings) error {
	err := r.db.QueryRow(
		`UPDATE users
		SET app_settings = $1
		WHERE id = $2`,
		settings_model.AppSettings,
		user_id,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) GetCountUserOwnedChats(user_id int) (count int, err error) {
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM chats 
		WHERE owner_id = $1`,
		user_id,
	).Scan(&count)
	return
}

func (r *UsersRepo) UserExistsByID(user_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = $1
		)`,
		user_id,
	).Scan(&exists)

	return
}

func (r *UsersRepo) UserExistsByDomain(user_domain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM units
			INNER JOIN users
			ON users.id = units.id
			WHERE units.domain = $1
		)`,
		user_domain,
	).Scan(&exists)

	return
}

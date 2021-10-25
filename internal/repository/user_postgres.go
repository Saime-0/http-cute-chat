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

func (r *UsersRepo) GetUserIdByInput(input_model models.UserInput) (id int, err error) {
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

func (r *UsersRepo) GetUserInfoByDomain(domain string) (user models.UserInfo, err error) {
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

func (r *UsersRepo) GetUserInfoByID(id int) (user models.UserInfo, err error) {
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

func (r *UsersRepo) IsUserExistsByInput(input_model models.UserInput) bool {
	is_exists := false
	err := r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM units 
			INNER JOIN users 
			ON units.id = users.id 
			WHERE units.domain = $1 AND units.name = $2
			)`,
		input_model.Domain,
		input_model.Name,
	).Scan(&is_exists)
	if err != nil || !is_exists {
		return false
	}
	return true
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

func (r *UsersRepo) CreateNewUserRefreshSession(user_id int, session_model *models.RefreshSession) (sessions_count int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO refresh_sessions (user_id, refresh_token, user_agent, exp, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		user_id,
		session_model.RefreshToken,
		session_model.UserAgent,
		session_model.Exp,
		session_model.CreatedAt,
	).Err()
	if err != nil {
		return
	}
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM refresh_sessions
		WHERE user_id = $1`,
		user_id,
	).Scan(&sessions_count)
	if err != nil {
		return
	}
	// ? todo:
	// INSERT INTO scientist (id, firstname) VALUES (3, 'chel');
	// UPDATE scientist SET counter = counter + 1 WHERE ID = 3;
	// DELETE FROM scientist WHERE counter = 3 and id = 3;
	// or: вернуть counter после создания сессии и если сессий больше 5 то DeleteOldestSession(user_id)
	return
}
func (r *UsersRepo) DeleteOldestSession(user_id int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM refresh_sessions 
		WHERE ctid IN(SELECT ctid FROM refresh_sessions WHERE user_id=$1 LIMIT 1)`,
		user_id,
	).Err()
	if err != nil {
		return
	}
	return
}
func (r *UsersRepo) FindSessionByComparedToken(token string) (session_id int, user_id int, err error) {
	err = r.db.QueryRow(
		`SELECT id, user_id
		FROM refresh_sessions
		WHERE refresh_token = $1`,
		token,
	).Scan(
		&session_id,
		&user_id,
	)
	if err != nil {
		return
	}
	return

}
func (r *UsersRepo) UpdateRefreshSession(session_id int, session_model *models.RefreshSession) (err error) {
	err = r.db.QueryRow(
		`UPDATE refresh_sessions
		SET refresh_token = $2, user_agent = $3, exp = $4, created_at = $5
		WHERE id = $1`,
		session_id,
		session_model.RefreshToken,
		session_model.UserAgent,
		session_model.Exp,
		session_model.CreatedAt,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

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

func (r *UsersRepo) CreateUser(u *models.CreateUser) (id int, err error) {
	err = r.db.QueryRow(
		`WITH u AS (
			INSERT INTO units (domain, name) 
			VALUES ($1, $2) 
			RETURNING id
			) 
		INSERT INTO users (id, app_settings) 
		SELECT u.id FROM u 
		RETURNING id`,
		u.Domain,
		u.Name,
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

func (r *UsersRepo) GetUserIdByInput(input models.UserInput) (id int, err error) {
	err = r.db.QueryRow(
		`SELECT units.id
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1 AND users.name = $2`,
		input.Domain,
		input.Name,
	).Scan(&id)
	if err != nil {
		return
	}
	return
}

func (r *UsersRepo) GetUserInfoByDomain(domain string) (user models.UserInfo, err error) {
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
	)
	if err != nil {
		return // user, err
	}
	return
}

func (r *UsersRepo) GetUserInfoByID(id int) (user models.UserInfo, err error) {
	err = r.db.QueryRow(
		`SELECT units.id,units.domain,units.name,users.app_settings 
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

func (r *UsersRepo) GetListUsersByName(name string) (chats models.ListUserInfo, err error) {

	return
}

func (r *UsersRepo) GetListChatsOwnedUser(user_id int) (chats models.ListChatInfo, err error) {

	return
}

func (r *UsersRepo) GetListChatsUser(user_id int) (chats models.ListChatInfo, err error) {

	return
}

func (r *UsersRepo) IsUserExistsByInput(input models.UserInput) bool {
	id := -1
	err := r.db.QueryRow(
		`SELECT units.id
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.domain = $1 AND users.name = $2`,
		input.Domain,
		input.Name,
	).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func (r *UsersRepo) GetUserSettings(user_id int) (settings *models.UserSettings, err error) {
	err = r.db.QueryRow(
		`SELECT units.app_settings
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

func (r *UsersRepo) UpdateUserData(user_id int, inp *models.UpdateUserData) error {
	if inp.Domain != "" {
		err := r.db.QueryRow(
			`UPDATE units
			SET domain = $1
			WHERE id = $2`,
			inp.Domain,
			user_id,
		).Err()
		if err != nil {
			return err
		}
	}
	if inp.Name != "" {
		err := r.db.QueryRow(
			`UPDATE units
			SET name = $1
			WHERE id = $2`,
			inp.Name,
			user_id,
		).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UsersRepo) UpdateUserSettings(user_id int, inp *models.UpdateUserSettings) error {
	err := r.db.QueryRow(
		`UPDATE users
		SET app_settings = $1
		WHERE id = $2`,
		inp.AppSettings,
		user_id,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) CreateNewUserRefreshSession(user_id int, s *models.RefreshSession) (sessions_count int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO refresh_sessions (user_id, refresh_token, user_agent, exp, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		user_id,
		s.RefreshToken,
		s.UserAgent,
		s.Exp,
		s.CreatedAt,
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
func (r *UsersRepo) UpdateRefreshSession(session_id int, s *models.RefreshSession) (err error) {
	err = r.db.QueryRow(
		`UPDATE refresh_sessions
		SET refresh_token = $2, user_agent = $3, exp = $4, created_at = $5
		WHERE id = $1`,
		session_id,
		s.RefreshToken,
		s.UserAgent,
		s.Exp,
		s.CreatedAt,
	).Err()
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) CreateNewUserRefreshSession(user_id int, session_model *models.RefreshSession) (sessions_count int, err error) {
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

	return
}
func (r *AuthRepo) DeleteOldestSession(user_id int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM refresh_sessions 
		WHERE ctid IN(SELECT ctid FROM refresh_sessions WHERE user_id=$1 LIMIT 1)`,
		user_id,
	).Err()

	return
}
func (r *AuthRepo) FindSessionByComparedToken(token string) (session_id int, user_id int, err error) {
	err = r.db.QueryRow(
		`SELECT id, user_id
		FROM refresh_sessions
		WHERE refresh_token = $1`,
		token,
	).Scan(
		&session_id,
		&user_id,
	)

	return
}
func (r *AuthRepo) UpdateRefreshSession(session_id int, session_model *models.RefreshSession) (err error) {
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

	return nil
}

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

func (r *AuthRepo) CreateNewUserRefreshSession(userId int, sessionModel *models.RefreshSession) (sessionsCount int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO refresh_sessions (user_id, refresh_token, user_agent, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userId,
		sessionModel.RefreshToken,
		sessionModel.UserAgent,
		sessionModel.Exp,
		sessionModel.CreatedAt,
	).Err()
	if err != nil {
		return
	}
	err = r.db.QueryRow(
		`SELECT count(*)
		FROM refresh_sessions
		WHERE user_id = $1`,
		userId,
	).Scan(&sessionsCount)

	return
}
func (r *AuthRepo) DeleteOldestSession(userId int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM refresh_sessions 
		WHERE ctid IN(SELECT ctid FROM refresh_sessions WHERE user_id=$1 LIMIT 1)`,
		userId,
	).Err()

	return
}
func (r *AuthRepo) FindSessionByComparedToken(token string) (sessionId int, userId int, err error) {
	err = r.db.QueryRow(
		`SELECT id, user_id
		FROM refresh_sessions
		WHERE refresh_token = $1`,
		token,
	).Scan(
		&sessionId,
		&userId,
	)

	return
}
func (r *AuthRepo) UpdateRefreshSession(sessionId int, sessionModel *models.RefreshSession) (err error) {
	err = r.db.QueryRow(
		`UPDATE refresh_sessions
		SET refresh_token = $2, user_agent = $3, expires_at = $4, created_at = $5
		WHERE id = $1`,
		sessionId,
		sessionModel.RefreshToken,
		sessionModel.UserAgent,
		sessionModel.Exp,
		sessionModel.CreatedAt,
	).Err()

	return nil
}

package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/rules"

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

func (r *AuthRepo) CreateRefreshSession(userId int, sessionModel *models.RefreshSession, overflowDelete bool) (expiresAt int64, err error) {
	err = r.db.QueryRow(`
		INSERT INTO refresh_sessions (user_id, refresh_token, user_agent, expires_at)
		VALUES ($1, $2, $3, unix_utc_now($4))
		RETURNING expires_at`,
		userId,
		sessionModel.RefreshToken,
		sessionModel.UserAgent,
		sessionModel.Lifetime,
	).Scan(&expiresAt)
	if err != nil {
		println(err.Error()) // debug
		return
	}
	if overflowDelete {
		err = r.OverflowDelete(userId, rules.MaxRefreshSession)
		if err != nil {
			println(err.Error()) // debug
			return
		}
	}

	return
}
func (r *AuthRepo) OverflowDelete(userId, limit int) (err error) {
	err = r.db.QueryRow(
		`DELETE FROM refresh_sessions 
		WHERE id IN(                 
		    WITH session_count AS (
		        SELECT count(1) AS val
				FROM refresh_sessions
				WHERE user_id = $1
		        GROUP BY user_id
		    )
		    SELECT id
		    FROM refresh_sessions 
		    WHERE coalesce((select val from session_count) > $2, false) = true AND user_id = $1
		    LIMIT abs((select val from session_count) - $2)
		    )`,
		userId,
		limit,
	).Err()
	if err != nil {
		println(err.Error()) // debug
	}
	return
}

func (r *AuthRepo) FindSessionByComparedToken(token string) (sessionId int, userId int, err error) {
	err = r.db.QueryRow(`
		SELECT id, user_id
		FROM refresh_sessions
		WHERE refresh_token = $1`,
		token,
	).Scan(
		&sessionId,
		&userId,
	)

	return
}

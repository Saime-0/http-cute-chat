package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/models"
)

type PreparesRepo struct {
	db *sql.DB
}

func NewPreparesRepo(db *sql.DB) *PreparesRepo {
	return &PreparesRepo{
		db: db,
	}
}

func (r *PreparesRepo) ScheduleInvites(before int64) ([]*models.ScheduleInvite, error) {
	var invites []*models.ScheduleInvite
	rows, err := r.db.Query(`
		SELECT code, chat_id, expires_at
		FROM invites
		WHERE $1 = 0 OR expires_at <= $1
		`,
		before,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		inv := &models.ScheduleInvite{}
		if err = rows.Scan(&inv.Code, &inv.ChatID, &inv.Exp); err != nil {
			return nil, err
		}
		invites = append(invites, inv)
	}

	return invites, nil
}

func (r *PreparesRepo) ScheduleRegisterSessions(before int64) ([]*models.ScheduleRegisterSession, error) {
	var sessions []*models.ScheduleRegisterSession

	rows, err := r.db.Query(`
		SELECT email, expires_at
		FROM registration_session
		WHERE $1 = 0 OR expires_at <= $1
		`,
		before,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		rs := &models.ScheduleRegisterSession{}
		if err = rows.Scan(&rs.Email, &rs.Exp); err != nil {
			return nil, err
		}
		sessions = append(sessions, rs)
	}

	return sessions, nil
}

func (r *PreparesRepo) ScheduleRefreshSessions(before int64) ([]*models.ScheduleRefreshSession, error) {
	var sessions []*models.ScheduleRefreshSession

	rows, err := r.db.Query(`
		SELECT id, expires_at
		FROM refresh_sessions
		WHERE $1 = 0 OR expires_at <= $1
		`,
		before,
	)
	if err != nil {
		return sessions, err
	}
	defer rows.Close()
	for rows.Next() {
		rs := &models.ScheduleRefreshSession{}
		if err = rows.Scan(&rs.ID, &rs.Exp); err != nil {
			return sessions, err
		}
		sessions = append(sessions, rs)
	}

	return sessions, nil
}

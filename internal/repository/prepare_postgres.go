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

func (r *PreparesRepo) ScheduleInvites() ([]*models.ScheduleInvite, error) {
	var invites []*models.ScheduleInvite

	rows, err := r.db.Query(`
		SELECT code, chat_id, expires_at
		FROM invites`,
	)
	if err != nil {
		println("ScanInvites:", err.Error()) // debug
		return invites, err
	}
	defer rows.Close()
	for rows.Next() {
		inv := &models.ScheduleInvite{}
		if err = rows.Scan(&inv.Code, &inv.ChatID, &inv.Exp); err != nil {
			println("ScheduleInvites:", err.Error()) // debug
			return invites, err
		}
		invites = append(invites, inv)
	}

	return invites, nil
}

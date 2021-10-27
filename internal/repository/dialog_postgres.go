package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type DialogsRepo struct {
	db *sql.DB
}

func NewDialogsRepo(db *sql.DB) *DialogsRepo {
	return &DialogsRepo{
		db: db,
	}
}

func (r *DialogsRepo) GetDialogIDBetweenUsers(user1_id int, user2_id int) (dialog_id int, err error) {
	err = r.db.QueryRow(
		`SELECT id
		FROM dialogs
		WHERE user1 = $1 AND user2 = $2 OR user1 = $2 AND user2 = $1`,
		user1_id,
		user2_id,
	).Scan(&dialog_id)
	if err != nil {
		return
	}
	return
}

func (r *DialogsRepo) GetCompanions(user_id int) (users models.ListUserInfo, err error) {
	rows, err := r.db.Query(
		`SELECT units.id,units.domain,units.name 
		FROM units INNER JOIN users 
		ON units.id = users.id 
		WHERE units.id IN (
			SELECT user2
			FROM dialogs
			WHERE user1 = $1
			UNION
			SELECT user1
			FROM dialogs
			WHERE user2 = $1
		)`,
		user_id,
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

func (r *DialogsRepo) DialogExistsBetweenUsers(user1_id int, user2_id int) (exits bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM dialogs 
			WHERE user1 = $1 AND user2 = $2 OR user1 = $2 AND user2 = $1
			)`,
		user1_id,
		user2_id,
	).Scan(&exits)
	return
}
func (r *DialogsRepo) CreateDialogBetweenUser(user1_id int, user2_id int) (dialog_id int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO dialogs (user1, user2)
		VALUES ($1, $2)
		RETURNING id`,
		user1_id,
		user2_id,
	).Scan(&dialog_id)
	if err != nil {
		return
	}
	return
}

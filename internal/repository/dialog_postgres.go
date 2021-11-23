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

func (r *DialogsRepo) GetDialogIDBetweenUsers(user1Id int, user2Id int) (dialogId int, err error) {
	err = r.db.QueryRow(
		`SELECT id
		FROM dialogs
		WHERE user1 = $1 AND user2 = $2 OR user1 = $2 AND user2 = $1`,
		user1Id,
		user2Id,
	).Scan(&dialogId)
	if err != nil {
		return
	}
	return
}

func (r *DialogsRepo) GetCompanions(userId int) (users models.ListUserInfo, err error) {
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
		userId,
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

func (r *DialogsRepo) DialogExistsBetweenUsers(user1Id int, user2Id int) (exits bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM dialogs 
			WHERE user1 = $1 AND user2 = $2 OR user1 = $2 AND user2 = $1
			)`,
		user1Id,
		user2Id,
	).Scan(&exits)
	return
}
func (r *DialogsRepo) CreateDialogBetweenUser(user1Id int, user2Id int) (dialogId int, err error) {
	err = r.db.QueryRow(
		`INSERT INTO dialogs (user1, user2)
		VALUES ($1, $2)
		RETURNING id`,
		user1Id,
		user2Id,
	).Scan(&dialogId)
	if err != nil {
		return
	}
	return
}

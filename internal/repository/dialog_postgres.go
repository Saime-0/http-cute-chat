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

func (r *DialogsRepo) CreateMessage(dialog_id int, message_model *models.CreateMessage) (message_id int, err error) {
	err = r.db.QueryRow( // todo: add record to "dialogs" if not exists
		`WITH m AS (
			INSERT INTO messages (reply_to, author, body, type)
			VALUES($2, $3, $4, $5)
			RETURNING id
			)
		INSERT INTO dialog_msg_pool (dialog_id, message_id) 
		SELECT $1, m.id
		FROM m`,
		message_model.ReplyTo,
		message_model.Author,
		message_model.Body,
		message_model.Type,
	).Scan(&message_id)
	if err != nil {
		return
	}
	return
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
		`SELECT units.id,units.domain,units.name,users.app_settings 
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

func (r *DialogsRepo) GetMessages(dialog_id int) (messages models.MessagesList, err error) {
	rows, err := r.db.Query(
		`SELECT id, reply_to, author, body, type
		FROM messages
		WHERE id IN (
			SELECT message_id 
			FROM dialog_msg_pool 
			WHERE dialog_id = $1 
			LIMIT 20
			)`,
		dialog_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.MessageInfo{}
		if err = rows.Scan(&m.ID, &m.ReplyTo, &m.Author, &m.Body, &m.Type); err != nil {
			return
		}
		messages.Messages = append(messages.Messages, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *DialogsRepo) GetMessageInfo(message_id int, dialog_id int) (message models.MessageInfo, err error) {
	err = r.db.QueryRow(
		`SELECT messages.id, messages.reply_to, messages.author, messages.body, messages.type
		FROM messages
		INNER JOIN dialog_msg_pool
		ON messages.id = dialog_msg_pool.message_id
		WHERE dialog_id = $1 AND message_id = $2`,
		dialog_id,
		message_id,
	).Scan(
		&message.ID,
		&message.ReplyTo,
		&message.Author,
		&message.Body,
		&message.Type,
	)
	if err != nil {
		return
	}
	return
}

package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type MessagesRepo struct {
	db *sql.DB
}

func NewMessagesRepo(db *sql.DB) *MessagesRepo {
	return &MessagesRepo{
		db: db,
	}
}

func (r *MessagesRepo) MessageExistsByID(message_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE id = $1
			)`,
		message_id,
	).Scan(&exists)

	return
}

func (r *MessagesRepo) MessageAvailableOnDialog(message_id int, dialog_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM messages
			INNER JOIN dialog_msg_pool
			ON messages.id = dialog_msg_pool.message_id
			WHERE dialog_id = $1 AND message_id = $2
			)`,
		dialog_id,
		message_id,
	).Scan(&exists)

	return
}

func (r *MessagesRepo) MessageAvailableOnRoom(message_id int, room_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM messages
			INNER JOIN room_msg_pool
			ON messages.id = room_msg_pool.message_id
			WHERE room_id = $1 AND message_id = $2
			)`,
		room_id,
		message_id,
	).Scan(&exists)

	return
}

func (r *MessagesRepo) GetMessageFromDialog(message_id int, dialog_id int) (message models.MessageInfo, err error) {
	err = r.db.QueryRow(
		`SELECT messages.id, COALESCE(messages.reply_to, 0), messages.author, messages.body, messages.type
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

	return
}

func (r *MessagesRepo) GetMessageFromRoom(message_id int, room_id int) (message models.MessageInfo, err error) {
	err = r.db.QueryRow(
		`SELECT messages.id, COALESCE(messages.reply_to, 0), messages.author, messages.body, messages.type
		FROM messages
		INNER JOIN room_msg_pool
		ON messages.id = room_msg_pool.message_id
		WHERE room_id = $1 AND message_id = $2`,
		room_id,
		message_id,
	).Scan(
		&message.ID,
		&message.ReplyTo,
		&message.Author,
		&message.Body,
		&message.Type,
	)

	return
}

func (r *MessagesRepo) GetMessagesFromDialog(dialog_id int, offset int) (messages models.MessagesList, err error) {
	rows, err := r.db.Query(
		`SELECT id, COALESCE(reply_to, 0), author, body, type
		FROM messages
		WHERE id IN (
			SELECT message_id 
			FROM dialog_msg_pool 
			WHERE dialog_id = $1 
			LIMIT 20
			OFFSET $2
			)`,
		dialog_id,
		offset,
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
	// if !rows.NextResultSet() {
	// 	return
	// }
	return
}

func (r *MessagesRepo) GetMessagesFromRoom(room_id int, created_after int, offset int) (messages models.MessagesList, err error) {
	rows, err := r.db.Query(
		`SELECT id, COALESCE(reply_to, 0), author, body, type
		FROM messages
		WHERE id IN (
			SELECT message_id 
			FROM room_msg_pool 
			WHERE room_id = $1  AND time > $3
			LIMIT 20
			OFFSET $2
			)`,
		room_id,
		offset,
		created_after,
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

	return
}

func (r *MessagesRepo) CreateMessageInDialog(dialog_id int, message_model *models.CreateMessage) (message_id int, err error) {
	err = r.db.QueryRow(
		`WITH m AS (
			INSERT INTO messages (reply_to, author, body, type)
			VALUES (NULLIF($2, 0), $3, $4, $5)
			RETURNING id
		)
		INSERT INTO dialog_msg_pool (dialog_id, message_id) 
		SELECT $1, m.id
		FROM m
		RETURNING message_id`,
		dialog_id,
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

func (r *MessagesRepo) CreateMessageInRoom(room_id int, message_model *models.CreateMessage) (message_id int, err error) {
	err = r.db.QueryRow(
		`WITH m AS (
			INSERT INTO messages (reply_to, author, body, type)
			VALUES (NULLIF($2, 0), $3, $4, $5)
			RETURNING id
		)
		INSERT INTO room_msg_pool (room_id, message_id) 
		SELECT $1, m.id
		FROM m
		RETURNING message_id`,
		room_id,
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

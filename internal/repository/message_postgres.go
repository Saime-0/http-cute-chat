package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
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

func (r *MessagesRepo) MessageExistsByID(messageId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE id = $1
			)`,
		messageId,
	).Scan(&exists)

	return
}

func (r *MessagesRepo) MessageAvailableOnRoom(messageId int, roomId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM messages
			WHERE id = $1 AND room_id = $2
			)`,
		messageId,
		roomId,
	).Scan(&exists)

	return
}

func (r *MessagesRepo) GetMessageFromDialog(messageId int, dialogId int) (message models.MessageInfo, err error) {
	err = r.db.QueryRow(
		`SELECT messages.id, COALESCE(messages.reply_to, 0), messages.author, messages.body, messages.type
		FROM messages
		INNER JOIN dialog_msg_pool
		ON messages.id = dialog_msg_pool.message_id
		WHERE dialog_id = $1 AND message_id = $2`,
		dialogId,
		messageId,
	).Scan(
		&message.ID,
		&message.ReplyTo,
		&message.Author,
		&message.Body,
		&message.Type,
	)

	return
}

func (r *MessagesRepo) Message(messageId int) (*model.Message, error) {
	message := &model.Message{}
	err := r.db.QueryRow(
		`SELECT id, body, type, created_at
		FROM messages
		WHERE id = $1`,
		messageId,
	).Scan(
		&message.ID,
		&message.Body,
		&message.Type,
		&message.Date,
	)

	return message, err
}

func (r *MessagesRepo) GetMessagesFromDialog(dialogId int, offset int) (messages models.MessagesList, err error) {
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
		dialogId,
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

func (r *MessagesRepo) GetMessagesFromRoom(roomId int, createdAfter int, offset int) (messages models.MessagesList, err error) {
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
		roomId,
		offset,
		createdAfter,
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

func (r *MessagesRepo) CreateMessageInDialog(dialogId int, messageModel *models.CreateMessage) (messageId int, err error) {
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
		dialogId,
		messageModel.ReplyTo,
		messageModel.Author,
		messageModel.Body,
		rules.UserMsg,
	).Scan(&messageId)
	if err != nil {
		return
	}
	return
}

func (r *MessagesRepo) CreateMessageInRoom(roomId int, msgType rules.MessageType, messageModel *models.CreateMessage) (messageId int, err error) {
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
		roomId,
		messageModel.ReplyTo,
		messageModel.Author,
		messageModel.Body,
		msgType,
	).Scan(&messageId)
	if err != nil {
		return
	}
	return
}

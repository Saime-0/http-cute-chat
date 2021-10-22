package repository

import (
	"database/sql"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type RoomsRepo struct {
	db *sql.DB
}

func NewRoomsRepoo(db *sql.DB) *RoomsRepo {
	return &RoomsRepo{
		db: db,
	}
}

func (r *RoomsRepo) CreateMessage(room_id int, message_model *models.CreateMessage) (message_id int, err error) {
	err = r.db.QueryRow(
		`WITH m AS (
			INSERT INTO messages (reply_to, author, body, type)
			VALUES($2, $3, $4, $5)
			RETURNING id
			)
		INSERT INTO room_msg_pool (room_id, message_id) 
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

func (r *RoomsRepo) GetMessages(room_id int) (messages models.MessagesList, err error) {
	rows, err := r.db.Query(
		`SELECT id, reply_to, author, body, type
		FROM messages
		WHERE id IN (
			SELECT message_id 
			FROM room_msg_pool 
			WHERE room_id = $1 
			LIMIT 20
			)`,
		room_id,
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

func (r *RoomsRepo) IsRoomExistsByID(room_id int) (is_exists bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM rooms WHERE id=$1)`,
		room_id,
	).Scan(&is_exists)
	if err != nil || !is_exists {
		return
	}
	return
}

func (r *RoomsRepo) CreateRoom(room_model *models.CreateRoom) (room_id int, err error) {
	err = r.db.QueryRow(
		`INSER INTO rooms (chat_id, parent_room, name, desc)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		room_model.ChatID,
		room_model.ParentRoom,
		room_model.Name,
		room_model.Desc,
	).Scan(&room_id)
	if err != nil {
		return
	}
	return
}

func (r *RoomsRepo) UpdateRoomData(room_id int, room_model *models.UpdateRoomData) (err error) {
	if room_model.Name != "" {
		err = r.db.QueryRow(
			`UPDATE rooms
			SET name = $2
			WHERE id = $1`,
			room_id,
			room_model.Name,
		).Err()
		if err != nil {
			return
		}
	}
	if room_model.Desc != "" {
		err = r.db.QueryRow(
			`UPDATE rooms
			SET desc = $2
			WHERE id = $1`,
			room_id,
			room_model.Desc,
		).Err()
		if err != nil {
			return
		}
	}
	return
}

func (r *RoomsRepo) GetChatIDByRoomID(room_id int) (chat_id int, err error) {
	err = r.db.QueryRow(
		`SELECT chat_id
		FROM rooms
		WHERE id = $1`,
		room_id,
	).Err()
	if err != nil {
		return
	}
	return
}

func (r *RoomsRepo) GetMessageInfo(message_id int, room_id int) (message models.MessageInfo, err error) {
	err = r.db.QueryRow(
		`SELECT messages.id, messages.reply_to, messages.author, messages.body, messages.type
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
	if err != nil {
		return
	}
	return
}

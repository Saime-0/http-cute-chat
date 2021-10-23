package repository

import (
	"database/sql"
	"errors"

	"github.com/saime-0/http-cute-chat/internal/models"
)

type RoomsRepo struct {
	db *sql.DB
}

func NewRoomsRepo(db *sql.DB) *RoomsRepo {
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
	if room_model.ParentRoom != 0 {
		parent_room_is_child := false
		err = r.db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM room_relation WHERE child_id=$1)`,
			room_model.ParentRoom,
		).Scan(&parent_room_is_child)
		if err != nil {
			return
		}
		if parent_room_is_child {
			return 0, errors.New("parent's room with a child")
		}
	}
	err = r.db.QueryRow(
		`INSERT INTO rooms (chat_id, name, note)
		VALUES ($1, $2, $3)
		RETURNING id`,
		room_model.ChatID,
		room_model.Name,
		room_model.Note,
	).Scan(&room_id)
	if err != nil {
		return
	}
	if room_model.ParentRoom != 0 {
		err = r.db.QueryRow(
			`INSERT INTO room_relation (parent_id, child_id)
			VALUES ($1, $2)`,
			room_model.ParentRoom,
			room_id,
		).Err()
		if err != nil {
			return
		}
	}
	return
}

// ? почему не работает?
func (r *RoomsRepo) GetChatRooms(chat_id int) (rooms models.ListRoomInfo, err error) {
	rows, err := r.db.Query(
		`SELECT rooms.id, room_relation.parent_id, rooms.name, rooms.note
		FROM rooms
		INNER JOIN room_relation
		ON rooms.id = room_relation.child_id
		WHERE rooms.chat_id = $1`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.RoomInfo{}
		if err = rows.Scan(&m.ID, &m.ParentRoom, &m.Name, &m.Note); err != nil {
			return
		}
		rooms.Rooms = append(rooms.Rooms, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *RoomsRepo) GetRoomInfo(room_id int) (room models.RoomInfo, err error) {
	err = r.db.QueryRow(
		`SELECT rooms.id, room_relation.parent_id, rooms.name, rooms.note
		FROM rooms
		INNER JOIN room_relation
		ON rooms.id = room_relation.child_id
		WHERE rooms.id = $1`,
		room_id,
	).Scan()
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
	if room_model.Note != "" {
		err = r.db.QueryRow(
			`UPDATE rooms
			SET note = $2
			WHERE id = $1`,
			room_id,
			room_model.Note,
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

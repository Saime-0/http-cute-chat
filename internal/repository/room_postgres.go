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

func (r *RoomsRepo) RoomExistsByID(room_id int) (is_exists bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM rooms WHERE id=$1)`,
		room_id,
	).Scan(&is_exists)
	if err != nil || !is_exists {
		return
	}
	return
}

func (r *RoomsRepo) CreateRoom(chat_id int, room_model *models.CreateRoom) (room_id int, err error) {
	if room_model.ParentID != 0 {
		parent_room_is_child := false
		err = r.db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM rooms WHERE parent_id IS NOT NULL)`,
		).Scan(&parent_room_is_child)
		if err != nil {
			return
		}
		if parent_room_is_child {
			return 0, errors.New("parent's room with a child")
		}
	}
	err = r.db.QueryRow(
		`INSERT INTO rooms (chat_id, parent_id, name, note, private)
		VALUES ($1, NULLIF($2, 0), $3, $4, $5)
		RETURNING id`,
		chat_id,
		room_model.ParentID,
		room_model.Name,
		room_model.Note,
		room_model.Private,
	).Scan(&room_id)
	if err != nil {
		return
	}
	return
}

func (r *RoomsRepo) GetChatRooms(chat_id int) (rooms models.ListRoomInfo, err error) {
	rows, err := r.db.Query(
		`SELECT id, COALESCE(parent_id, 0), name, note, private
		FROM rooms
		WHERE chat_id = $1`,
		chat_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := models.RoomInfo{}
		if err = rows.Scan(&m.ID, &m.ParentID, &m.Name, &m.Note, &m.Private); err != nil {
			return
		}
		rooms.Rooms = append(rooms.Rooms, m)
	}
	if !rows.NextResultSet() {
		return
	}
	return
}

func (r *RoomsRepo) GetRoom(room_id int) (room models.RoomInfo, err error) {
	err = r.db.QueryRow(
		`SELECT id, COALESCE(parent_id, 0), name, note, private
		FROM rooms
		WHERE id = $1`,
		room_id,
	).Scan(
		&room.ID,
		&room.ParentID,
		&room.Name,
		&room.Note,
		&room.Private,
	)

	return
}

func (r *RoomsRepo) UpdateRoomData(room_id int, room_model *models.UpdateRoomData) (err error) {
	if room_model.Name != "" {
		err = r.db.QueryRow(
			`UPDATE rooms
			SET name = $2, private = $3
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
			SET note = $2, private = $3
			WHERE id = $1`,
			room_id,
			room_model.Note,
			room_model.Private,
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
	).Scan(&chat_id)
	if err != nil {
		return
	}
	return
}

func (r *RoomsRepo) RoomIsPrivate(room_id int) (private bool) {
	r.db.QueryRow(
		`SELECT private
		FROM rooms
		WHERE id = $1`,
		room_id,
	).Scan(&private)

	return
}

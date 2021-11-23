package repository

import (
	"database/sql"
	"encoding/json"
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

func (r *RoomsRepo) RoomExistsByID(roomId int) (isExists bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM rooms WHERE id=$1)`,
		roomId,
	).Scan(&isExists)
	if err != nil || !isExists {
		return
	}
	return
}

func (r *RoomsRepo) CreateRoom(chatId int, roomModel *models.CreateRoom) (roomId int, err error) {
	if roomModel.ParentID != 0 {
		parentRoomIsChild := false
		err = r.db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM rooms WHERE parent_id IS NOT NULL)`,
		).Scan(&parentRoomIsChild)
		if err != nil {
			return
		}
		if parentRoomIsChild {
			return 0, errors.New("parent's room with a child")
		}
	}
	err = r.db.QueryRow(
		`INSERT INTO rooms (chat_id, parent_id, name, note, private)
		VALUES ($1, NULLIF($2, 0), $3, $4, $5)
		RETURNING id`,
		chatId,
		roomModel.ParentID,
		roomModel.Name,
		roomModel.Note,
		roomModel.Private,
	).Scan(&roomId)
	if err != nil {
		return
	}
	return
}

func (r *RoomsRepo) GetChatRooms(chatId int) (rooms models.ListRoomInfo, err error) {
	rows, err := r.db.Query(
		`SELECT id, COALESCE(parent_id, 0), name, note, private
		FROM rooms
		WHERE chat_id = $1`,
		chatId,
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

func (r *RoomsRepo) GetRoom(roomId int) (room models.RoomInfo, err error) {
	err = r.db.QueryRow(
		`SELECT id, COALESCE(parent_id, 0), name, note, private
		FROM rooms
		WHERE id = $1`,
		roomId,
	).Scan(
		&room.ID,
		&room.ParentID,
		&room.Name,
		&room.Note,
		&room.Private,
	)

	return
}

func (r *RoomsRepo) UpdateRoomData(roomId int, roomModel *models.UpdateRoomData) (err error) {
	if roomModel.Name != "" {
		err = r.db.QueryRow(
			`UPDATE rooms
			SET name = $2, private = $3
			WHERE id = $1`,
			roomId,
			roomModel.Name,
		).Err()
		if err != nil {
			return
		}
	}
	if roomModel.Note != "" {
		err = r.db.QueryRow(
			`UPDATE rooms
			SET note = $2, private = $3
			WHERE id = $1`,
			roomId,
			roomModel.Note,
			roomModel.Private,
		).Err()
		if err != nil {
			return
		}
	}

	return
}

func (r *RoomsRepo) GetChatIDByRoomID(roomId int) (chatId int, err error) {
	err = r.db.QueryRow(
		`SELECT chat_id
		FROM rooms
		WHERE id = $1`,
		roomId,
	).Scan(&chatId)
	if err != nil {
		return
	}
	return
}

func (r *RoomsRepo) RoomIsPrivate(roomId int) (private bool) {
	r.db.QueryRow(
		`SELECT private
		FROM rooms
		WHERE id = $1`,
		roomId,
	).Scan(&private)

	return
}

func (r *RoomsRepo) GetRoomForm(roomId int) (form models.FormPattern, err error) {
	var format string
	err = r.db.QueryRow(
		`SELECT COALESCE(msg_format, '')
		FROM rooms
		WHERE room_id = $1`,
		roomId,
	).Scan(&format)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(format), &form)

	return
}

func (r *RoomsRepo) UpdateRoomForm(roomId int, format string) (err error) {
	err = r.db.QueryRow(
		`UPDATE rooms
		SET msg_format = NULLIF($2, '')
		WHERE id = $1`,
		roomId,
		format,
	).Err()

	return
}

func (r *RoomsRepo) RoomFormIsSet(roomId int) (isSet bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1
			FROM rooms
			WHERE id = $1 AND msg_format IS NOT NULL
		)`,
		roomId,
	).Scan(&isSet)

	return
}

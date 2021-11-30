package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"strconv"
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

func (r *RoomsRepo) CreateRoom(chatId int, input *model.CreateRoomInput) (roomId int, err error) {
	var format *string
	if input.MsgFormat != nil {
		marshal, err := json.Marshal(*input.MsgFormat)
		if err != nil {
			return
		}
		*format = string(marshal)
	}
	err = r.db.QueryRow(
		`INSERT INTO rooms (chat_id, parent_id, name, note, msg_format)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		chatId,
		input.Parent,
		input.Name,
		input.Note,
		input.MsgFormat,
	).Scan(&roomId)
	if err != nil {
		return
	}
	var allows *string
	if input.Restricts != nil {
		var allowsDb []models.AllowsDB
		if input.Restricts.AllowWrite != nil {
			for _, role := range input.Restricts.AllowWrite.Roles {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowWrite,
					Group:  rules.AllowRoles,
					Value:  strconv.Itoa(role),
				})
			}
			for _, char := range input.Restricts.AllowWrite.Chars {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowWrite,
					Group:  rules.AllowChars,
					Value:  char.String(),
				})
			}
			for _, member := range input.Restricts.AllowWrite.Members {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowWrite,
					Group:  rules.AllowChars,
					Value:  strconv.Itoa(member),
				})
			}
		}

		if input.Restricts.AllowRead != nil {
			for _, role := range input.Restricts.AllowRead.Roles {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowRead,
					Group:  rules.AllowRoles,
					Value:  strconv.Itoa(role),
				})
			}
			for _, char := range input.Restricts.AllowRead.Chars {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowRead,
					Group:  rules.AllowChars,
					Value:  char.String(),
				})
			}
			for _, member := range input.Restricts.AllowRead.Members {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowRead,
					Group:  rules.AllowChars,
					Value:  strconv.Itoa(member),
				})
			}
		}

		if len(allowsDb) != 0 {
			for _, allow := range allowsDb {
				*allows += fmt.Sprintf(",(%d, '%s','%s','%s')", roomId, allow.Action, allow.Group, allow.Value)
			}
			*allows = kit.TrimFirstRune(*allows)
		}

	}

	if allows != nil {

		_, err = r.db.Exec(`INSERT INTO allows 
    		(room_id, action_type, group_type, value)	
			VALUES` + *allows)
		if err != nil {
			return
		}

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
		WHERE id = $1`,
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

func (r *RoomsRepo) GetAllows(room_id int) (allows models.Allows, err error) {
	rows, err := r.db.Query(
		`SELECT action_type, group_type, value
	FROM allows
	WHERE room_id = $1`,
		room_id)
	if err != nil {
		return
	}
	defer rows.Close()

	allows = models.Allows{
		Read:  models.AllowHolders{},
		Write: models.AllowHolders{},
	}
	for rows.Next() {
		d := models.AllowsDB{}
		if err = rows.Scan(&d.Action, &d.Group, &d.Value); err != nil {
			return
		}

		matchAction := func() *models.AllowHolders {
			switch d.Action {
			case rules.AllowRead:
				return &allows.Read
			case rules.AllowWrite:
				return &allows.Write
			default:
				panic("GetAllows lose action matching")
			}
		}
		h := matchAction()
		switch d.Group {
		case rules.AllowChars:
			switch d.Value {
			case string(rules.Admin):
				h.Chars = append(h.Chars, rules.Admin)
			case string(rules.Moder):
				h.Chars = append(h.Chars, rules.Moder)
			default:
				panic("GetAllows not identify value type")
			}
		case rules.AllowRoles:
			value, err := strconv.Atoi(d.Value)
			if err != nil {
				panic("GetAllows int was expected but an error was found")
			}
			h.Roles = append(h.Roles, value)
		case rules.AllowUsers:
			value, err := strconv.Atoi(d.Value)
			if err != nil {
				panic("GetAllows int was expected but an error was found")
			}
			h.Users = append(h.Users, value)
		default:
			panic("GetAllows not identify group type")
		}

	}

	return
}

func (r *RoomsRepo) HasParent(roomId int) (has bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM rooms 
			WHERE id = $1 AND parent_id IS NOT NULL 
    	)`,
		roomId,
	).Scan(&has)

	return
}

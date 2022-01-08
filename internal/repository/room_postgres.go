package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/tlog"
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

func (r *RoomsRepo) CreateRoom(inp *model.CreateRoomInput) (*model.CreateRoom, error) {
	var (
		roomId int
		err    error
		room   = &model.CreateRoom{}
		format *string
	)

	if inp.Form != nil {
		var marshal []byte
		marshal, err = json.Marshal(*inp.Form)
		if err != nil {
			println("CreateRoom:", err.Error()) // debug
			return room, err
		}
		*format = string(marshal)
	}

	err = r.db.QueryRow(`
		INSERT INTO rooms (chat_id, parent_id, name, note, msg_format)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, chat_id, name, parent_id, note`,
		inp.ChatID,
		inp.Parent,
		inp.Name,
		inp.Note,
		inp.Form,
	).Scan(
		&room.ID,
		&room.ChatID,
		&room.Name,
		&room.ParentID,
		&room.Note,
	)
	if err != nil {
		println("CreateRoom:", err.Error()) // debug
		return room, err
	}

	if inp.Allows != nil {
		var allows string
		for _, allow := range inp.Allows.Allows {
			allows += fmt.Sprintf(",(%d, '%s','%s','%s')", roomId, allow.Action, allow.Group, allow.Value)
		}
		allows = kit.TrimFirstRune(allows)

		//language=PostgreSQL
		err = r.db.QueryRow(`
			INSERT INTO allows (room_id, action_type, group_type, value)
			VALUES` + allows,
		).Err()
		if err != nil {
			println("CreateRoom:", err.Error()) // debug
			return room, err
		}

	}
	return room, nil
}

func (r *RoomsRepo) Room(roomId int) (*model.Room, error) {
	room := &model.Room{
		Chat: &model.Chat{
			Unit: &model.Unit{},
		},
	}
	err := r.db.QueryRow(
		`SELECT  id, chat_id, parent_id, name, note
		FROM rooms
		WHERE id = $1`,
		roomId,
	).Scan(
		&room.RoomID,
		&room.Chat.Unit.ID,
		&room.ParentID,
		&room.Name,
		&room.Note,
	)

	return room, err
}

func (r *RoomsRepo) UpdateRoom(roomId int, inp *model.UpdateRoomInput) (*model.UpdateRoom, error) {
	room := &model.UpdateRoom{}
	err := r.db.QueryRow(`
		UPDATE rooms
		SET 
		    name = COALESCE($2::VARCHAR, name), 
		    parent_id = COALESCE($3::BIGINT,parent_id), 
		    note = COALESCE($4::VARCHAR,note)
		WHERE id = $1
		RETURNING id, name, parent_id, note`,
		roomId,
		inp.Name,
		inp.ParentID,
		inp.Note,
	).Scan(
		&room.ID,
		&room.Name,
		&room.ParentID,
		&room.Note,
	)
	if err != nil {
		println("UpdateRoom:", err.Error()) // debug
	}

	return room, err
}

func (r *RoomsRepo) DeleteRoom(roomID int) (*model.DeleteRoom, error) {
	room := &model.DeleteRoom{}
	err := r.db.QueryRow(`
		DELETE FROM rooms
		WHERE id = $1
		RETURNING id`,
		roomID,
	).Scan(&room.ID)
	if err != nil {
		println("DeleteRoom:", err.Error()) // debug
	}

	return room, err
}

func (r *RoomsRepo) DeleteAllow(allowID int) (*model.DeleteAllow, error) {
	allow := &model.DeleteAllow{}
	err := r.db.QueryRow(`
		DELETE FROM allows
		WHERE id = $1
		RETURNING id`,
		allowID,
	).Scan(&allow.AllowID)
	if err != nil {
		println("DeleteAllow:", err.Error()) // debug
	}

	return allow, err
}
func (r *RoomsRepo) AllowExists(roomID int, inp *model.AllowInput) (exists bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM allows
			WHERE room_id = $1 
			  AND action_type = $2 
			  AND group_type = $3 
			  AND value = $4
		)
		`,
		roomID,
		inp.Action,
		inp.Group,
		inp.Value,
	).Scan(&exists)
	if err != nil {
		println("AllowExists:", err.Error()) // debug
	}

	return
}

func (r *RoomsRepo) CreateAllow(roomID int, inp *model.AllowInput) (*model.CreateAllow, error) {
	allow := &model.CreateAllow{
		Allow: &model.Allow{},
	}
	err := r.db.QueryRow(`
		INSERT INTO allows (room_id, action_type, group_type, value)
		VALUES ($1, $2, $3, $4)
		RETURNING id, room_id, action_type, group_type, value`,
		roomID,
		inp.Action,
		inp.Group,
		inp.Value,
	).Scan(
		&allow.RoomID,
		&allow.Allow.ID,
		&allow.Allow.Action,
		&allow.Allow.Group,
		&allow.Allow.Value,
	)
	if err != nil {
		println("CreateAllow:", err.Error()) // debug
	}

	return allow, err
}

func (r *RoomsRepo) GetChatIDByRoomID(roomId int) (chatId int, err error) {
	err = r.db.QueryRow(
		`SELECT chat_id
		FROM rooms
		WHERE id = $1`,
		roomId,
	).Scan(&chatId)
	if err != nil {
		println("GetChatIDByRoomID:", err.Error()) // debug
	}
	return
}
func (r *RoomsRepo) GetChatIDByAllowID(allowID int) (chatId int, err error) {
	err = r.db.QueryRow(`
		SELECT chat_id
		FROM rooms r
		JOIN allows a on r.id = a.room_id
		WHERE a.id = $1`,
		allowID,
	).Scan(&chatId)
	if err != nil {
		println("GetChatIDByAllowID:", err.Error()) // debug
	}
	return
}

func (r *RoomsRepo) RoomForm(roomId int) *model.Form {
	var err error
	tl := tlog.Start("RoomsRepo > RoomForm [rid:" + strconv.Itoa(roomId) + "]")
	defer tl.Fine()
	var (
		format *string
		form   *model.Form
	)
	err = r.db.QueryRow(
		`SELECT msg_format
		FROM rooms
		WHERE id = $1`,
		roomId,
	).Scan(&format)
	if err != nil {
		println("RoomForm:", err.Error()) // debug
		return nil
	}
	if format == nil {
		return nil
	}

	err = json.Unmarshal([]byte(*format), &form)
	if err != nil {
		println("RoomForm:", err.Error()) // debug
		return nil
	}

	return form
}

func (r *RoomsRepo) UpdateRoomForm(roomId int, form *string) (err error) {
	err = r.db.QueryRow(
		`UPDATE rooms
		SET msg_format = $2
		WHERE id = $1`,
		roomId,
		form,
	).Err()

	return
}

func (r *RoomsRepo) FormIsSet(roomId int) (have bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(
				SELECT 1
				FROM rooms
				WHERE id = $1 AND msg_format IS NOT NULL
			)`,
		roomId,
	).Scan(&have)
	if err != nil {
		println("FormIsSet:", err.Error())
	}

	return
}

func (r *RoomsRepo) Allows(roomID int) (*model.Allows, error) {
	allows := &model.Allows{
		Room: &model.Room{
			RoomID: roomID,
		},
		Allows: []*model.Allow{},
	}
	rows, err := r.db.Query(`
		SELECT action_type, group_type, value
		FROM allows
		WHERE room_id = $1`,
		roomID,
	)
	if err != nil {
		println("Allows:", err.Error())
		return allows, err
	}
	defer rows.Close()

	for rows.Next() {
		allow := &model.Allow{}
		if err = rows.Scan(&allow.Action, &allow.Group, &allow.Value); err != nil {
			panic(err)
		}
		allows.Allows = append(allows.Allows, allow)
	}
	return allows, nil
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

func (r *RoomsRepo) Allowed(action rules.AllowActionType, roomId int, holder *models.AllowHolder) (yes bool) {
	tl := tlog.Start("RoomsRepo > Allowed [rid:" + strconv.Itoa(roomId) + ",uid:" + strconv.Itoa(holder.UserID) + "]")
	defer tl.Fine()
	err := r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1 
		    FROM chats
		    JOIN rooms r on chats.id = r.chat_id
		    LEFT JOIN allows a on r.id = a.room_id
		    WHERE r.id = $2 
		    	AND (
			        action_type IS NULL 
			        OR action_type = $1 
			            AND (
							group_type = 'ROLE' AND value = $3::VARCHAR 
							OR group_type = 'CHAR' AND value = $4::VARCHAR  
							OR group_type = 'MEMBER' AND value = $6::VARCHAR
						)
			        OR owner_id = $5::BIGINT 
		    	)
	    )`,
		action,
		roomId,
		holder.RoleID,
		holder.Char,
		holder.UserID,
		holder.MemberID,
	).Scan(&yes)
	if err != nil {
		println("Allowed:", err.Error()) //debug
	}

	return
}

func (r *RoomsRepo) AllowHolder(userId, chatId int) (*models.AllowHolder, error) {
	tl := tlog.Start("RoomsRepo > AllowHolder [uid:" + strconv.Itoa(userId) + ",cid:" + strconv.Itoa(chatId) + "]")
	defer tl.Fine()
	holder := &models.AllowHolder{
		RoleID: nil,
		Char:   "",
		UserID: 0,
	}
	err := r.db.QueryRow(
		`SELECT id, role_id, char
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(
		&holder.MemberID,
		&holder.RoleID,
		&holder.Char,
	)
	if err != nil {
		println("AllowsHolder:", err.Error()) //debug
		return nil, err
	}
	holder.UserID = userId

	return holder, err
}

func (r *RoomsRepo) AllowsIsSet(roomId int) (have bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
    		SELECT 1 
    		FROM allows
    		WHERE room_id = $1
		)`,
		roomId,
	).Scan(&have)

	return
}

func (r *RoomsRepo) FindRooms(inp *model.FindRooms, params *model.Params) *model.Rooms {
	rooms := &model.Rooms{
		Rooms: []*model.Room{},
	}
	if inp.NameFragment != nil {
		*inp.NameFragment = "%" + *inp.NameFragment + "%"
	}
	// language=PostgreSQL
	rows, err := r.db.Query(`
		SELECT rooms.id, chats.id,  name, parent_id, note
		FROM rooms
		JOIN chats ON rooms.chat_id = chats.id
		WHERE chat_id = $1
			AND (
			    $2::BIGINT IS NULL 
			    OR rooms.id = $2 
			)
			AND (
			    $3::VARCHAR IS NULL 
			    OR rooms.name ILIKE $3
			)
			AND (
			    $5::fetch_type IN (NULL, 'NEUTRAL') 
			    OR $5::fetch_type = 'POSITIVE' 
			           AND parent_id IS NOT NULL AND (
			               $4::BIGINT IS NULL 
			               OR parent_id = $4
			           )
			    OR $5::fetch_type = 'NEGATIVE' 
			           AND parent_id IS NULL
			)
		LIMIT $6
		OFFSET $7
	`,
		inp.ChatID,
		inp.RoomID,
		inp.NameFragment,
		inp.ParentID,
		inp.IsChild,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		println("FindRooms:", err.Error()) // debug
		return rooms
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Room{
			Chat: &model.Chat{
				Unit: &model.Unit{},
			},
		}
		if err = rows.Scan(&m.RoomID, &m.Chat.Unit.ID, &m.Name, &m.ParentID, &m.Note); err != nil {
			println("rows.scan:", err.Error()) // debug
			return rooms
		}

		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms
}

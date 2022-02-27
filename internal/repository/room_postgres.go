package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

type RoomsRepo struct {
	db *sql.DB
}

func NewRoomsRepo(db *sql.DB) *RoomsRepo {
	return &RoomsRepo{
		db: db,
	}
}

func (r *RoomsRepo) RoomExistsByID(roomID int) (isExists bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS (
    		SELECT 1 
    		FROM rooms 
    		WHERE id=$1
		)
		`,
		roomID,
	).Scan(&isExists)
	if err != nil || !isExists {
		return
	}
	return
}

func (r *RoomsRepo) CreateRoom(inp *model.CreateRoomInput) (*model.CreateRoom, error) {
	var (
		err    error
		room   = &model.CreateRoom{}
		format *string
	)

	if inp.Form != nil {
		var marshal []byte
		marshal, err = json.Marshal(*inp.Form)
		if err != nil {
			return nil, err
		}
		format = kit.StringPtr(string(marshal))
	}

	err = r.db.QueryRow(`
		INSERT INTO rooms (chat_id, parent_id, name, note, msg_format)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, chat_id, name, parent_id, note`,
		inp.ChatID,
		inp.Parent,
		inp.Name,
		inp.Note,
		//inp.Form,
		format,
	).Scan(
		&room.ID,
		&room.ChatID,
		&room.Name,
		&room.ParentID,
		&room.Note,
	)
	if err != nil {
		return nil, err
	}

	if inp.Allows != nil {
		var allows string
		for _, allow := range inp.Allows.Allows {
			allows += fmt.Sprintf(",(%d, '%s','%s','%s')", room.ID, allow.Action, allow.Group, allow.Value)
		}
		allows = kit.TrimFirstRune(allows)

		//language=PostgreSQL
		err = r.db.QueryRow(`
			INSERT INTO allows (room_id, action_type, group_type, value)
			VALUES` + allows,
		).Err()
		if err != nil {
			return nil, err
		}

	}
	return room, nil
}

func (r *RoomsRepo) Room(roomID int) (*model.Room, error) {
	room := &model.Room{
		Chat: &model.Chat{
			Unit: new(model.Unit),
		},
	}
	err := r.db.QueryRow(`
		SELECT  id, chat_id, parent_id, name, note
		FROM rooms
		WHERE id = $1
		`,
		roomID,
	).Scan(
		&room.RoomID,
		&room.Chat.Unit.ID,
		&room.ParentID,
		&room.Name,
		&room.Note,
	)

	return room, err
}

func (r *RoomsRepo) UpdateRoom(roomID int, inp *model.UpdateRoomInput) (*model.UpdateRoom, error) {
	room := &model.UpdateRoom{}
	err := r.db.QueryRow(`
		UPDATE rooms
		SET 
		    name = COALESCE($2::VARCHAR, name), 
		    parent_id = COALESCE($3::BIGINT,parent_id), 
		    note = COALESCE($4::VARCHAR,note)
		WHERE id = $1
		RETURNING id, name, parent_id, note`,
		roomID,
		inp.Name,
		inp.ParentID,
		inp.Note,
	).Scan(
		&room.ID,
		&room.Name,
		&room.ParentID,
		&room.Note,
	)

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

	return allow, err
}

func (r *RoomsRepo) AllowsExists(roomID int, allows *model.AllowsInput) (desired bool, err error) {
	err = r.db.QueryRow(`
		SELECT *
		FROM unnest($2::findallow[]) elem (act,gr,val)
		LEFT JOIN allows a ON a.room_id =  $1
	        AND elem.act = a.action_type::VARCHAR
	        AND elem.gr = a.group_type::VARCHAR
	        AND elem.val = a.value
		`,
		roomID,
		pq.Array(allows.Allows),
	).Scan(&desired)

	return
}

func (r *RoomsRepo) CreateAllows(roomID int, inp *model.AllowsInput) (*model.CreateAllows, error) {
	allows := &model.CreateAllows{
		RoomID: roomID,
		Allows: []*model.Allow{},
	}

	sqlArr := ""
	for _, v := range inp.Allows {

		sqlArr += fmt.Sprintf(",(%d, '%s', '%s', '%s')", roomID, v.Action, v.Group, v.Value)
	}
	sqlArr = kit.TrimFirstRune(sqlArr)
	rows, err := r.db.Query(`
		INSERT INTO allows (room_id, action_type, group_type, value)
		VALUES ` + sqlArr + `
		RETURNING id, action_type, group_type, value`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		allow := &model.Allow{}
		if err = rows.Scan(&allow.ID, &allow.Action, &allow.Group, &allow.Value); err != nil {
			panic(err)
		}
		allows.Allows = append(allows.Allows, allow)
	}

	return allows, err
}

func (r *RoomsRepo) GetChatIDByRoomID(roomID int) (chatID int, err error) {
	err = r.db.QueryRow(
		`SELECT chat_id
		FROM rooms
		WHERE id = $1`,
		roomID,
	).Scan(&chatID)

	return
}
func (r *RoomsRepo) GetChatIDByAllowID(allowID int) (chatID int, err error) {
	err = r.db.QueryRow(`
		SELECT coalesce((
		    SELECT chat_id
			FROM rooms r
			JOIN allows a on r.id = a.room_id
			WHERE a.id = $1
		), 0)`,
		allowID,
	).Scan(&chatID)

	return
}

func (r *RoomsRepo) RoomForm(roomID int) (*model.Form, error) {
	var (
		format *string
		form   *model.Form
	)
	err := r.db.QueryRow(
		`SELECT msg_format
		FROM rooms
		WHERE id = $1`,
		roomID,
	).Scan(&format)
	if err != nil {
		return nil, err
	}
	if format == nil {
		return nil, nil
	}

	err = json.Unmarshal([]byte(*format), &form)
	if err != nil {
		return nil, err
	}

	return form, err
}

func (r *RoomsRepo) UpdateRoomForm(roomID int, form *string) (err error) {
	err = r.db.QueryRow(`
		UPDATE rooms
		SET msg_format = $2
		WHERE id = $1`,
		roomID,
		form,
	).Err()

	return
}

func (r *RoomsRepo) FormIsSet(roomID int) (have bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM rooms
			WHERE id = $1 AND msg_format IS NOT NULL
		)`,
		roomID,
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

func (r *RoomsRepo) HasParent(roomID int) (has bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM rooms
			WHERE id = $1 AND parent_id IS NOT NULL 
		)`,
		roomID,
	).Scan(&has)

	return
}

func (r *RoomsRepo) Allowed(action model.ActionType, roomID int, holder *models.AllowHolder) (yes bool) {
	err := r.db.QueryRow(`
		SELECT EXISTS(
		    SELECT 1 
		    FROM chats
		    JOIN rooms r on chats.id = r.chat_id
		    LEFT JOIN (
		        SELECT *
		        FROM allows
		        WHERE room_id = $2 AND action_type = $1
		    ) a on r.id = a.room_id
		    WHERE r.id = $2 
				AND (
			        action_type IS NULL 
			        OR (
						group_type = 'ROLE' AND value = $3::VARCHAR 
						OR group_type = 'CHAR' AND value = $4::VARCHAR  
						OR group_type = 'MEMBER' AND value = $6::VARCHAR
					)
			        OR owner_id = $5::BIGINT 
		    	)
	    )`,
		action,
		roomID,
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

func (r *RoomsRepo) AllowHolder(userID, chatID int) (*models.AllowHolder, error) {

	holder := &models.AllowHolder{}
	err := r.db.QueryRow(
		`SELECT id, role_id, char
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userID,
		chatID,
	).Scan(
		&holder.MemberID,
		&holder.RoleID,
		&holder.Char,
	)
	if err != nil {
		println("AllowsHolder:", err.Error()) //debug
		return nil, err
	}
	holder.UserID = userID

	return holder, err
}

func (r *RoomsRepo) AllowsIsSet(roomID int) (have bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
    		SELECT 1 
    		FROM allows
    		WHERE room_id = $1
		)`,
		roomID,
	).Scan(&have)

	return
}

func (r *RoomsRepo) FindRooms(inp *model.FindRooms, params *model.Params) (*model.Rooms, error) {
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
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Room{
			Chat: &model.Chat{
				Unit: new(model.Unit),
			},
		}
		if err = rows.Scan(&m.RoomID, &m.Chat.Unit.ID, &m.Name, &m.ParentID, &m.Note); err != nil {
			return nil, err
		}

		rooms.Rooms = append(rooms.Rooms, m)
	}

	return rooms, nil
}

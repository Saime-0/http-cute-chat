package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
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

func (r *RoomsRepo) CreateRoom(chatId int, input *model.CreateRoomInput) (roomId int, err error) {
	var format *string
	if input.Form != nil {
		var marshal []byte
		marshal, err = json.Marshal(*input.Form)
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
		input.Form,
	).Scan(&roomId)
	if err != nil {
		return
	}
	var allows *string
	if input.Allows != nil {
		var allowsDb []models.AllowsDB
		if input.Allows.AllowWrite != nil {
			for _, role := range input.Allows.AllowWrite.Roles {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowWrite,
					Group:  rules.AllowRoles,
					Value:  strconv.Itoa(role),
				})
			}
			for _, char := range input.Allows.AllowWrite.Chars {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowWrite,
					Group:  rules.AllowChars,
					Value:  char.String(),
				})
			}
			for _, member := range input.Allows.AllowWrite.Members {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowWrite,
					Group:  rules.AllowChars,
					Value:  strconv.Itoa(member),
				})
			}
		}

		if input.Allows.AllowRead != nil {
			for _, role := range input.Allows.AllowRead.Roles {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowRead,
					Group:  rules.AllowRoles,
					Value:  strconv.Itoa(role),
				})
			}
			for _, char := range input.Allows.AllowRead.Chars {
				allowsDb = append(allowsDb, models.AllowsDB{
					Action: rules.AllowRead,
					Group:  rules.AllowChars,
					Value:  char.String(),
				})
			}
			for _, member := range input.Allows.AllowRead.Members {
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

func (r *RoomsRepo) GetAllows(roomId int) (*model.Allows, error) {
	rows, err := r.db.Query(
		`SELECT action_type, group_type, value
	FROM allows
	WHERE room_id = $1`,
		roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	_allows := &models.Allows{
		Read:  models.AllowHolders{},
		Write: models.AllowHolders{},
	}
	for rows.Next() {
		d := models.AllowsDB{}
		if err = rows.Scan(&d.Action, &d.Group, &d.Value); err != nil {
			return nil, err
		}
		var h *models.AllowHolders
		switch d.Action {
		case rules.AllowRead:
			h = &_allows.Read
		case rules.AllowWrite:
			h = &_allows.Write
		default:
			panic("GetAllows lose action matching")
		}
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

	allows := &model.Allows{
		Room: nil, // outside
		AllowRead: &model.PermissionHolders{
			Roles: &model.Roles{
				Roles: []*model.Role{},
			},
			Chars: &model.Chars{
				Chars: []model.CharType{},
			},
			Members: &model.Members{
				Members: []*model.Member{},
			},
		},
		AllowWrite: &model.PermissionHolders{
			Roles: &model.Roles{
				Roles: []*model.Role{},
			},
			Chars: &model.Chars{
				Chars: []model.CharType{},
			},
			Members: &model.Members{
				Members: []*model.Member{},
			},
		},
	}
	chatId, err := r.GetChatIDByRoomID(roomId)
	if err != nil {
		return nil, err
	}
	configAllows := func(aholdres *models.AllowHolders, phold *model.PermissionHolders) error {
		if len(aholdres.Roles) != 0 {
			roles, err := r.RolesByArray(&aholdres.Roles)
			if err != nil {
				return err
			}
			phold.Roles = roles
		}
		if len(aholdres.Users) != 0 {
			members, err := r.MembersByArray(chatId, &aholdres.Users)
			if err != nil {
				return err
			}
			phold.Members = members

		}
		if len(aholdres.Chars) != 0 {
			for _, char := range aholdres.Chars {
				phold.Chars.Chars = append(phold.Chars.Chars, model.CharType(char))
			}
		}
		return nil
	}
	if configAllows(&_allows.Read, allows.AllowRead) != nil ||
		configAllows(&_allows.Write, allows.AllowWrite) != nil {
		return nil, err
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
	err := r.db.QueryRow(
		`SELECT EXISTS(
	    SELECT 1 
	    FROM allows
	    WHERE 
	    	action_type = $1 AND 
			room_id = $2 AND
			(
				group_type = 'ROLES' AND value = $3 OR
				group_type = 'CHARS' AND value = $4 OR
				group_type = 'USERS' AND value = $5 
			)
	    )`,
		action,
		roomId,
		holder.RoleID,
		holder.Char,
		holder.UserID,
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
		`SELECT role_id, char
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`,
		userId,
		chatId,
	).Scan(
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

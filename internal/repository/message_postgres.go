package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"

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

func (r *MessagesRepo) Message(messageId int) (*model.Message, error) {
	message := &model.Message{
		Room: &model.Room{},
	}
	var (
		_replid *int
		_userID *int
	)
	err := r.db.QueryRow(
		`SELECT id, reply_to, user_id, room_id, body, type, created_at
		FROM messages
		WHERE id = $1`,
		messageId,
	).Scan(
		&message.ID,
		&_replid,
		&_userID,
		&message.Room.RoomID,
		&message.Body,
		&message.Type,
		&message.CreatedAt,
	)
	if _replid != nil {
		message.ReplyTo = &model.Message{
			ID: *_replid,
		}
	}
	if _userID != nil {
		message.User = &model.User{
			Unit: &model.Unit{ID: *_userID},
		}
	}
	return message, err
}

func (r *MessagesRepo) CreateMessageInRoom(inp *models.CreateMessage) (*model.NewMessage, error) {
	message := &model.NewMessage{}
	err := r.db.QueryRow(`
		INSERT INTO messages (room_id, reply_to, user_id, body, type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, room_id, reply_to, user_id, body, type, created_at
		`,
		inp.RoomID,
		inp.ReplyTo,
		inp.UserID,
		inp.Body,
		inp.Type,
	).Scan(
		&message.ID,
		&message.RoomID,
		&message.ReplyToID,
		&message.UserID,
		&message.Body,
		&message.MsgType,
		&message.CreatedAt,
	)
	if err != nil {
		println("CreateMessageInRoom:", err.Error()) // debug
		return message, err
	}
	return message, nil
}

func (r *MessagesRepo) MessagesFromRoom(roomId, chatId int, find *model.FindMessagesInRoom) *model.Messages {
	messages := &model.Messages{
		Messages: []*model.Message{},
	}

	var direction int8 = 1

	if find.Created == model.MessagesCreatedBefore {
		direction = -1
	}
	rows, err := r.db.Query(`
		SELECT id, reply_to, user_id, room_id, body, messages.type, created_at
		FROM messages
		WHERE room_id = $1
		    AND (
		        $2 = 0
		        OR (
		            $4 = 1 AND id >= $2 AND id <= ($2 + $3)
		            OR $4 = -1 AND id <= $2
		        )
		    )
		ORDER BY created_at DESC
		OFFSET 0
		LIMIT $3
		`,
		roomId,
		find.StartMessageID,
		find.Count,
		direction,
	)
	if err != nil {
		println("MessagesFromRoom:", err.Error()) // debug
		return messages
	}
	defer rows.Close()
	for rows.Next() {
		m := &model.Message{
			Room: &model.Room{
				Chat: &model.Chat{
					Unit: &model.Unit{ID: chatId},
				},
			},
		}
		var (
			_replid *int
			_userID *int
		)
		if err = rows.Scan(&m.ID, &_replid, &_userID, &m.Room.RoomID, &m.Body, &m.Type, &m.CreatedAt); err != nil {
			println("rows.scan:", err.Error()) // debug
			return messages
		}
		if _replid != nil {
			m.ReplyTo = &model.Message{
				ID: *_replid,
			}
		}
		if _userID != nil {
			m.User = &model.User{
				Unit: &model.Unit{ID: *_userID},
			}
		}
		messages.Messages = append(messages.Messages, m)
	}

	return messages
}

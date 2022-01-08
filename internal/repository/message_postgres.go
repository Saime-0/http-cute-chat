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
		_replid   *int
		_memberId *int
	)
	err := r.db.QueryRow(
		`SELECT id, reply_to, author, room_id, body, type, created_at
		FROM messages
		WHERE id = $1`,
		messageId,
	).Scan(
		&message.ID,
		&_replid,
		&_memberId,
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
	if _memberId != nil {
		message.Author = &model.Member{
			ID: *_memberId,
		}
	}
	return message, err
}

func (r *MessagesRepo) CreateMessageInRoom(inp *models.CreateMessage) (*model.NewMessage, error) {
	message := &model.NewMessage{}
	err := r.db.QueryRow(`
		INSERT INTO messages (reply_to, author, room_id, body, type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, room_id, reply_to, author, body, type, created_at
		`,
		inp.ReplyTo,
		inp.Author,
		inp.RoomID,
		inp.Body,
		inp.Type,
	).Scan(
		&message.ID,
		&message.RoomID,
		&message.ReplyToID,
		&message.AuthorID,
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

func (r *MessagesRepo) MessagesFromRoom(roomId, chatId int, find *model.FindMessagesInRoomByUnionInput, params *model.Params) *model.Messages {
	messages := &model.Messages{
		Messages: []*model.Message{},
	}
	rows, err := r.db.Query(`
		SELECT id, reply_to, author, room_id, body, messages.type, created_at
		FROM messages
		WHERE room_id = $1
		  AND (
		      $2::BIGINT IS NULL 
		      OR messages.created_at > $2
		  )
		  AND (
		      $3::BIGINT IS NULL 
		      OR messages.created_at <= $3
		  )
		ORDER BY created_at
		OFFSET $4 
		LIMIT $5
		`,
		roomId,
		find.AfterTime,
		find.BeforeTime,
		params.Offset,
		params.Limit,
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
			_replid   *int
			_memberId *int
		)
		if err = rows.Scan(&m.ID, &_replid, &_memberId, &m.Room.RoomID, &m.Body, &m.Type, &m.CreatedAt); err != nil {
			println("rows.scan:", err.Error()) // debug
			return messages
		}
		if _replid != nil {
			m.ReplyTo = &model.Message{
				ID: *_replid,
			}
		}
		if _memberId != nil {
			m.Author = &model.Member{
				ID: *_memberId,
			}
		}
		messages.Messages = append(messages.Messages, m)
	}

	return messages
}

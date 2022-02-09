package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *messageResult) isRequestResult() {}
func (r *messageInp) isRequestInput()     {}

type (
	messageInp struct {
		MessageID int
	}
	messageResult struct {
		Message *model.Message
	}
)

func (d *Dataloader) Message(messageID int) (*model.Message, error) {
	res := <-d.categories.Message.addBaseRequest(
		&messageInp{
			MessageID: messageID,
		},
		new(messageResult),
	)
	if res == nil {
		return nil, d.categories.Message.Error
	}
	return res.(*messageResult).Message, nil
}

func (c *parentCategory) message() {
	var (
		inp = c.Requests

		ptrs       []chanPtr
		messageIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		messageIDs = append(messageIDs, query.Inp.(*messageInp).MessageID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(id, 0), 
		       coalesce(reply_to, 0), 
		       coalesce(user_id, 0), 
		       coalesce(room_id, 0), 
		       coalesce(body, ''), 
		       coalesce(type, 'USER'),
		       coalesce(created_at, 0)
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, messageid)
		LEFT JOIN messages m ON m.id = inp.messageid
		`,
		pq.Array(ptrs),
		pq.Array(messageIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("message:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr     chanPtr
		replyTo *int
		userID  *int
	)
	for rows.Next() {
		m := &model.Message{Room: new(model.Room)}

		if err = rows.Scan(&ptr, &m.ID, &replyTo, &userID, &m.Room.RoomID, &m.Body, &m.Type, &m.CreatedAt); err != nil {
			//c.Dataloader.healer.Alert("message (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.ID == 0 {
			m = nil
			goto sendRequest
		}
		if replyTo != nil {
			m.ReplyTo = &model.Message{ID: *replyTo}
		}
		if userID != nil {
			m.User = &model.User{
				Unit: &model.Unit{ID: *userID},
			}
		}

	sendRequest:
		request := c.getRequest(ptr)
		request.Result.(*messageResult).Message = m
	}

	c.Error = nil
}

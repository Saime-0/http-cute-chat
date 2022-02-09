package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *roomResult) isRequestResult() {}
func (r *roomInp) isRequestInput()     {}

type (
	roomInp struct {
		RoomID int
	}
	roomResult struct {
		Room *model.Room
	}
)

func (d *Dataloader) Room(roomID int) (*model.Room, error) {
	res := <-d.categories.Room.addBaseRequest(
		&roomInp{
			RoomID: roomID,
		},
		new(roomResult),
	)
	if res == nil {
		return nil, d.categories.Room.Error
	}
	return res.(*roomResult).Room, nil
}

func (c *parentCategory) room() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		roomIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		roomIDs = append(roomIDs, query.Inp.(*roomInp).RoomID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(id, 0), 
		       coalesce(chat_id, 0), 
		       coalesce(parent_id, 0), 
		       coalesce(name, ''), 
		       coalesce(note, '')
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, roomid)
		LEFT JOIN rooms m ON m.id = inp.roomid
		`,
		pq.Array(ptrs),
		pq.Array(roomIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("room:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := &model.Room{
			Chat: &model.Chat{Unit: new(model.Unit)},
		}

		if err = rows.Scan(&ptr, &m.RoomID, &m.Chat.Unit.ID, &m.ParentID, &m.Name, &m.Note); err != nil {
			//c.Dataloader.healer.Alert("room (scan rows):" + err.Error())
			c.Error = err
			return
		}

		if m.RoomID == 0 {
			m = nil
		}
		request := c.getRequest(ptr)
		request.Result.(*roomResult).Room = m
	}

	c.Error = nil
}

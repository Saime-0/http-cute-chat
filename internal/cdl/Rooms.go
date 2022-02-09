package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *RoomsResult) isRequestResult() {}
func (r *RoomsInp) isRequestInput()     {}

type (
	RoomsResult struct {
		Rooms *model.Rooms
	}
	RoomsInp struct {
		ChatID int
	}
)

func (d *Dataloader) Rooms(chatID int) (*model.Rooms, error) {
	res := <-d.categories.Rooms.addBaseRequest(
		&RoomsInp{
			ChatID: chatID,
		},
		&RoomsResult{
			Rooms: &model.Rooms{
				Rooms: []*model.Room{},
			},
		},
	)
	if res == nil {
		return nil, d.categories.Rooms.Error
	}
	return res.(*RoomsResult).Rooms, nil
}

func (c *parentCategory) rooms() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		chatIDs []int
	)
	for _, query := range inp {
		chatIDs = append(chatIDs, query.Inp.(*RoomsInp).ChatID)
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr,
				coalesce(id, 0),
				coalesce(chat_id, 0),
				coalesce(parent_id, 0),
				coalesce(name, ''), 
				coalesce(note, '')
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, chatid)
		LEFT JOIN rooms c ON c.chat_id = inp.chatid
		`,
		pq.Array(ptrs),
		pq.Array(chatIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("rooms:" + err.Error())
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
			//c.Dataloader.healer.Alert("rooms (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.RoomID == 0 {
			continue
		}
		request := c.getRequest(ptr)
		request.Result.(*RoomsResult).Rooms.Rooms = append(request.Result.(*RoomsResult).Rooms.Rooms, m)
	}

	c.Error = nil
}

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
		SELECT arr.id, c.chat_id, c.id, parent_id, name, note
		FROM unnest($1::varchar[], $2::bigint[]) arr(id, chatid)
		JOIN rooms c ON c.chat_id = arr.chatid
		`,
		pq.Array(ptrs),
		pq.Array(chatIDs),
	)
	if err != nil {
		println("Rooms:", err.Error()) // debug
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		chatID int
		ptr    chanPtr
	)
	for rows.Next() {
		m := &model.Room{
			Chat: &model.Chat{
				Unit: &model.Unit{},
			},
		}

		if err = rows.Scan(&ptr, &chatID, &m.RoomID, &m.ParentID, &m.Name, &m.Note); err != nil {
			c.Error = err
			return
		}
		m.Chat.Unit.ID = chatID // для того чтобы в roomResolver.Chat можно было узнать ид чата который надо вернуть

		request, ok := c.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("c.Requests not exists by" + ptr)
		}
		request.Result.(*RoomsResult).Rooms.Rooms = append(request.Result.(*RoomsResult).Rooms.Rooms, m)
	}

	c.Error = nil
}

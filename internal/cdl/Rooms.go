package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

type RoomsResult struct {
	Rooms *model.Rooms
}

//func (r *RoomsResult) IsRequestResult() {}

type RoomsInp struct {
	ChatID int
}

func (r *RoomsCategory) AddRequest(chatID int) chan *RoomsResult {
	newClient := make(chan *RoomsResult)

	r.Requests[fmt.Sprint(newClient)] = &RoomsRequest{
		Ch: newClient,
		Inp: RoomsInp{
			ChatID: chatID,
		},
		Result: &RoomsResult{
			Rooms: &model.Rooms{
				Rooms: []*model.Room{},
			},
		},
	}

	defer r.OnAddRequest()
	return newClient
}

func (d *Dataloader) Rooms(chatID int) (*model.Rooms, error) {
	res := <-d.Categories.Rooms.AddRequest(chatID)
	if res == nil {
		println("реквест выполнился с ошибкой видимо") // debug
		return nil, d.Categories.Rooms.Error
	}
	println("реквест выполнился нормально") // debug
	return res.Rooms, nil
}

func (r *RoomsCategory) rooms() {
	var (
		inp = r.Requests

		ptrs    []chanPtr
		chatIDs []int
	)
	for _, query := range inp {
		chatIDs = append(chatIDs, query.Inp.ChatID)
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
	}

	rows, err := r.Dataloader.DB.Query(`
		SELECT arr.id, r.chat_id, r.id, parent_id, name, note
		FROM unnest($1::varchar[], $2::bigint[]) arr(id, chatid)
		JOIN rooms r ON r.chat_id = arr.chatid
		`,
		pq.Array(ptrs),
		pq.Array(chatIDs),
	)
	if err != nil {
		println("Rooms:", err.Error()) // debug
		//r.Requests = ?
		r.Error = err
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
			//r.Requests = ?
			r.Error = err
			return
		}
		m.Chat.Unit.ID = chatID // для того чтобы в roomResolver.Chat можно было узнать ид чата который надо вернуть

		request, ok := r.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("r.Requests not exists")
		}
		request.Result.Rooms.Rooms = append(request.Result.Rooms.Rooms, m)
	}

	r.Error = nil
}

type RoomsCategory struct {
	ParentCategory
	Requests map[chanPtr]*RoomsRequest
}

func (d *Dataloader) NewRoomsCategory() *RoomsCategory {
	c := &RoomsCategory{
		ParentCategory: ParentCategory{
			Dataloader:             d,
			RemainingRequestsCount: d.CapactiyRequests,
		},
		Requests: map[chanPtr]*RoomsRequest{},
	}
	c.LoadFn = func() {
		c.rooms()
		if c.Error != nil {
			for _, request := range c.Requests {
				select {
				case request.Ch <- nil:
				default:
				}
			}
		}
		for _, request := range c.Requests {
			select {
			case request.Ch <- request.Result:
			default:
			}
		}
	}
	c.PrepareForNextLaunch = func() {
		for ptr := range c.Requests {
			delete(c.Requests, ptr)
		}
	}
	return c
}

type RoomsRequest struct {
	Ch     chan *RoomsResult
	Inp    RoomsInp
	Result *RoomsResult
}

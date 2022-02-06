package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *RoomsResult) IsRequestResult() {}
func (r *RoomsInp) IsRequestInput()     {}

//func (r RoomsResultChan) IsRequestResultChan() {}

type (
	RoomsResult struct {
		Rooms *model.Rooms
	}
	RoomsInp struct {
		ChatID int
	}
	//RoomsResultChan chan *RoomsResult
)

func (c *ParentCategory) AddRoomsRequest(chatID int) BaseResultChan {
	newClient := make(BaseResultChan)
	c.Lock()
	c.Requests[fmt.Sprint(newClient)] = &BaseRequest{
		Ch: newClient,
		Inp: &RoomsInp{
			ChatID: chatID,
		},
		Result: &RoomsResult{
			Rooms: &model.Rooms{
				Rooms: []*model.Room{},
			},
		},
	}
	//c.Requests[fmt.Sprint(newClient)] = &RoomsRequest{
	//	Ch: newClient,
	//	Inp: RoomsInp{
	//		ChatID: chatID,
	//	},
	//	Result: &RoomsResult{
	//		Rooms: &model.Rooms{
	//			Rooms: []*model.Room{},
	//		},
	//	},
	//}
	c.Unlock()

	go c.OnAddRequest()
	return newClient
}

func (d *Dataloader) Rooms(chatID int) (*model.Rooms, error) {
	res := <-d.Categories.Rooms.AddRoomsRequest(chatID)
	if res == nil {
		println("реквест выполнился с ошибкой видимо") // debug
		return nil, d.Categories.Rooms.Error
	}
	println("реквест выполнился нормально") // debug
	return res.(*RoomsResult).Rooms, nil
}

func (c *ParentCategory) rooms() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		chatIDs []int
	)
	for _, query := range inp {
		chatIDs = append(chatIDs, query.Inp.(*RoomsInp).ChatID)
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
	}

	rows, err := c.Dataloader.DB.Query(`
		SELECT arr.id, c.chat_id, c.id, parent_id, name, note
		FROM unnest($1::varchar[], $2::bigint[]) arr(id, chatid)
		JOIN rooms c ON c.chat_id = arr.chatid
		`,
		pq.Array(ptrs),
		pq.Array(chatIDs),
	)
	if err != nil {
		println("Rooms:", err.Error()) // debug
		//c.Requests = ?
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
			//c.Requests = ?
			c.Error = err
			return
		}
		m.Chat.Unit.ID = chatID // для того чтобы в roomResolver.Chat можно было узнать ид чата который надо вернуть

		request, ok := c.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("c.Requests not exists")
		}
		request.Result.(*RoomsResult).Rooms.Rooms = append(request.Result.(*RoomsResult).Rooms.Rooms, m)
	}

	c.Error = nil
}

//type RoomsCategory struct {
//	ParentCategory
//	Requests map[chanPtr]*RoomsRequest
//}

func (d *Dataloader) NewRoomsCategory() *ParentCategory {
	c := d.NewParentCategory()
	c.LoadFn = c.rooms
	return c
}

//type RoomsRequest struct {
//	Ch     chan *RoomsResult
//	Inp    RoomsInp
//	Result *RoomsResult
//}

package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

type UserIsChatMemberInp struct {
	UserID int
	ChatID int
}

type UserIsChatMemberResult struct {
	Exists bool
}

func (r *UserIsChatMemberCategory) AddRequest(userID, chatID int) chan *UserIsChatMemberResult {
	newClient := make(chan *UserIsChatMemberResult)

	r.Requests[fmt.Sprint(newClient)] = &UserIsChatMemberRequest{
		Ch: newClient,
		Inp: UserIsChatMemberInp{
			UserID: userID,
			ChatID: chatID,
		},
		Result: &UserIsChatMemberResult{},
	}

	defer r.OnAddRequest()
	return newClient
}

func (d *Dataloader) UserIsChatMember(userID, chatID int) (bool, error) {
	res := <-d.Categories.UserIsChatMember.AddRequest(userID, chatID)
	if res == nil {
		println("реквест выполнился с ошибкой видимо") // debug
		return false, d.Categories.UserIsChatMember.Error
	}
	println("реквест выполнился нормально") // debug
	return res.Exists, nil
}

func (r *UserIsChatMemberCategory) userIsChatMember() {
	var (
		inp = r.Requests

		ptrs    []chanPtr
		userIDs []int
		chatIDs []int
	)
	for _, query := range inp {
		userIDs = append(userIDs, query.Inp.UserID)
		chatIDs = append(chatIDs, query.Inp.ChatID)
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
	}

	rows, err := r.Dataloader.DB.Query(`
		SELECT arr.id, m.id is not null 
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) arr(id, userid, chatid)
		LEFT JOIN chat_members m ON m.chat_id = arr.chatid AND m.user_id = arr.userid
		`,
		pq.Array(ptrs),
		pq.Array(userIDs),
		pq.Array(chatIDs),
	)
	if err != nil {
		println("userIsChatMember:", err.Error()) // debug
		//r.Requests = ?
		r.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		isMember bool
		ptr      chanPtr
	)
	for rows.Next() {

		if err = rows.Scan(&ptr, &isMember); err != nil {
			//r.Requests = ?
			r.Error = err
			return
		}

		request, ok := r.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("r.Requests not exists")
		}
		request.Result.Exists = isMember
	}

	r.Error = nil
}

type UserIsChatMemberCategory struct {
	ParentCategory
	Requests map[chanPtr]*UserIsChatMemberRequest
}

func (d *Dataloader) NewUserIsChatMemberCategory() *UserIsChatMemberCategory {
	c := &UserIsChatMemberCategory{
		ParentCategory: ParentCategory{
			Dataloader:             d,
			RemainingRequestsCount: d.CapactiyRequests,
		},
		Requests: map[chanPtr]*UserIsChatMemberRequest{},
	}
	c.LoadFn = func() {
		c.userIsChatMember()
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

type UserIsChatMemberRequest struct {
	Ch     chan *UserIsChatMemberResult
	Inp    UserIsChatMemberInp
	Result *UserIsChatMemberResult
}

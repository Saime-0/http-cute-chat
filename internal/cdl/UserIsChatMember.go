package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

func (r *UserIsChatMemberResult) IsRequestResult() {}
func (r *UserIsChatMemberInp) IsRequestInput()     {}

type (
	UserIsChatMemberInp struct {
		UserID int
		ChatID int
	}
	UserIsChatMemberResult struct {
		Exists bool
	}
)

func (c *ParentCategory) AddUserIsChatMemberRequest(userID, chatID int) BaseResultChan {

	newClient := make(BaseResultChan)
	c.Lock()
	c.Requests[fmt.Sprint(newClient)] = &BaseRequest{
		Ch: newClient,
		Inp: &UserIsChatMemberInp{
			UserID: userID,
			ChatID: chatID,
		},
		Result: &UserIsChatMemberResult{Exists: false},
	}
	c.Unlock()

	go c.OnAddRequest()
	return newClient
}

func (d *Dataloader) UserIsChatMember(userID, chatID int) (bool, error) {

	res := <-d.Categories.UserIsChatMember.AddUserIsChatMemberRequest(userID, chatID)
	if res == nil {
		return false, d.Categories.UserIsChatMember.Error
	}
	return res.(*UserIsChatMemberResult).Exists, nil
}

func (c *ParentCategory) userIsChatMember() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		userIDs []int
		chatIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		userIDs = append(userIDs, query.Inp.(*UserIsChatMemberInp).UserID)
		chatIDs = append(chatIDs, query.Inp.(*UserIsChatMemberInp).ChatID)
	}

	rows, err := c.Dataloader.DB.Query(`
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
		//c.Requests = ?
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		isMember bool
		ptr      chanPtr
	)
	for rows.Next() {

		if err = rows.Scan(&ptr, &isMember); err != nil {
			//c.Requests = ?
			c.Error = err
			return
		}

		request, ok := c.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("c.Requests not exists")
		}
		request.Result.(*UserIsChatMemberResult).Exists = isMember
	}

	c.Error = nil
}

func (d *Dataloader) NewUserIsChatMemberCategory() *ParentCategory {
	c := d.NewParentCategory()
	c.LoadFn = c.userIsChatMember
	return c
}

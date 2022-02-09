package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

func (r *UserIsChatMemberResult) isRequestResult() {}
func (r *UserIsChatMemberInp) isRequestInput()     {}

type (
	UserIsChatMemberInp struct {
		UserID int
		ChatID int
	}
	UserIsChatMemberResult struct {
		Exists bool
	}
)

func (d *Dataloader) UserIsChatMember(userID, chatID int) (bool, error) {

	res := <-d.categories.UserIsChatMember.addBaseRequest(
		&UserIsChatMemberInp{
			UserID: userID,
			ChatID: chatID,
		},
		new(UserIsChatMemberResult),
	)
	if res == nil {
		return false, d.categories.UserIsChatMember.Error
	}
	return res.(*UserIsChatMemberResult).Exists, nil
}

func (c *parentCategory) userIsChatMember() {
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

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, id is not null 
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) inp(ptr, userid, chatid)
		LEFT JOIN chat_members m ON m.chat_id = inp.chatid AND m.user_id = inp.userid
		`,
		pq.Array(ptrs),
		pq.Array(userIDs),
		pq.Array(chatIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("userIsChatMember:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr      chanPtr
		isMember bool
	)
	for rows.Next() {

		if err = rows.Scan(&ptr, &isMember); err != nil {
			//c.Dataloader.healer.Alert("userIsChatMember (scan rows):" + err.Error())
			c.Error = err
			return
		}

		request := c.getRequest(ptr)
		request.Result.(*UserIsChatMemberResult).Exists = isMember
	}

	c.Error = nil
}

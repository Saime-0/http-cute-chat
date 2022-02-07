package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

func (r *chatIDByMemberIDResult) isRequestResult() {}
func (r *chatIDByMemberIDInp) isRequestInput()     {}

type (
	chatIDByMemberIDInp struct {
		MemberID int
	}
	chatIDByMemberIDResult struct {
		ChatID int
	}
)

func (d *Dataloader) ChatIDByMemberID(memberID int) (int, error) {
	res := <-d.categories.ChatIDByMemberID.addBaseRequest(
		&chatIDByMemberIDInp{
			MemberID: memberID,
		},
		&chatIDByMemberIDResult{},
	)
	if res == nil {
		return 0, d.categories.ChatIDByMemberID.Error
	}
	return res.(*chatIDByMemberIDResult).ChatID, nil
}

func (c *parentCategory) chatIDByMemberID() {
	var (
		inp = c.Requests

		ptrs      []chanPtr
		memberIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		memberIDs = append(memberIDs, query.Inp.(*chatIDByMemberIDInp).MemberID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT arr.id, coalesce(m.chat_id, 0)
		FROM unnest($1::varchar[], $2::bigint[]) arr(id, memberid)
		LEFT JOIN chat_members m ON m.chat_id = arr.memberid = m.id
		`,
		pq.Array(ptrs),
		pq.Array(memberIDs),
	)
	if err != nil {
		println("chatIDByMemberID:", err.Error()) // debug
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr    chanPtr
		chatID int
	)
	for rows.Next() {
		if err = rows.Scan(&ptr, &chatID); err != nil {
			c.Error = err
			return
		}

		request, ok := c.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("c.Requests not exists")
		}
		request.Result.(*chatIDByMemberIDResult).ChatID = chatID
	}

	c.Error = nil
}
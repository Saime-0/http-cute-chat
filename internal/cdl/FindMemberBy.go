package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *findMemberByResult) isRequestResult() {}
func (r *findMemberByInp) isRequestInput()     {}

type (
	findMemberByInp struct {
		UserID int
		ChatID int
	}
	findMemberByResult struct {
		MemberID *int
	}
)

func (d *Dataloader) FindMemberBy(userID, chatID int) (*int, error) {
	res := <-d.categories.FindMemberBy.addBaseRequest(
		&findMemberByInp{
			ChatID: chatID,
			UserID: userID,
		},
		new(findMemberByResult),
	)
	if res == nil {
		return nil, d.categories.FindMemberBy.Error
	}
	return res.(*findMemberByResult).MemberID, nil
}

func (c *parentCategory) findMemberBy() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		userIDs []int
		chatIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		userIDs = append(userIDs, query.Inp.(*findMemberByInp).UserID)
		chatIDs = append(chatIDs, query.Inp.(*findMemberByInp).ChatID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, coalesce(id, 0)
		FROM unnest($1::varchar[], $2::bigint[], $3::bigint[]) inp(ptr, userid, chatid)
		LEFT JOIN chat_members m ON m.chat_id = inp.chatid AND m.user_id = inp.userid
		`,
		pq.Array(ptrs),
		pq.Array(userIDs),
		pq.Array(chatIDs),
	)
	if err != nil {
		c.Dataloader.healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr      chanPtr
		memberID *int
	)
	for rows.Next() {
		if err = rows.Scan(&ptr, &memberID); err != nil {
			c.Error = err
			return
		}

		request := c.getRequest(ptr)
		request.Result.(*findMemberByResult).MemberID = memberID
	}

	c.Error = nil
}

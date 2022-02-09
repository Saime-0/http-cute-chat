package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *memberRoleResult) isRequestResult() {}
func (r *memberRoleInp) isRequestInput()     {}

type (
	memberRoleInp struct {
		MemberID int
	}
	memberRoleResult struct {
		Role *model.Role
	}
)

func (d *Dataloader) MemberRole(memberID int) (*model.Role, error) {
	res := <-d.categories.MemberRole.addBaseRequest(
		&memberRoleInp{
			MemberID: memberID,
		},
		new(memberRoleResult),
	)
	if res == nil {
		return nil, d.categories.MemberRole.Error
	}
	return res.(*memberRoleResult).Role, nil
}

func (c *parentCategory) memberRole() {
	var (
		inp = c.Requests

		ptrs      []chanPtr
		memberIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		memberIDs = append(memberIDs, query.Inp.(*memberRoleInp).MemberID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, 
		       coalesce(r.id, 0),
		       coalesce(r.name, ''),
		       coalesce(r.color, '')
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, memberid)
		LEFT JOIN chat_members m ON m.id = inp.memberid
		LEFT JOIN  roles r ON m.role_id = r.id
		`,
		pq.Array(ptrs),
		pq.Array(memberIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("memberRole:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := new(model.Role)

		if err = rows.Scan(&ptr, &m.ID, &m.Name, &m.Color); err != nil {
			//c.Dataloader.healer.Alert("memberRole (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.ID == 0 {
			m = nil
		}

		request := c.getRequest(ptr)
		request.Result.(*memberRoleResult).Role = m
	}

	c.Error = nil
}

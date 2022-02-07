package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *UserResult) IsRequestResult() {}
func (r *UserInp) IsRequestInput()     {}

type (
	UserInp struct {
		UserID int
	}
	UserResult struct {
		User *model.User
	}
)

func (d *Dataloader) User(userID int) (*model.User, error) {
	res := <-d.Categories.User.AddBaseRequest(
		&UserInp{
			UserID: userID,
		},
		&UserResult{
			User: &model.User{
				Unit: &model.Unit{},
			},
		},
	)
	if res == nil {
		return nil, d.Categories.User.Error
	}
	return res.(*UserResult).User, nil
}

func (c *ParentCategory) user() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		userIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		userIDs = append(userIDs, query.Inp.(*UserInp).UserID)
	}

	rows, err := c.Dataloader.DB.Query(`
		SELECT arr.id, u.id, domain, name, type 
		FROM unnest($1::varchar[], $2::bigint[]) arr(id, userid)
		LEFT JOIN units u ON u.id = arr.userid AND u.type = 'USER'
		`,
		pq.Array(ptrs),
		pq.Array(userIDs),
	)
	if err != nil {
		println("user:", err.Error()) // debug
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := &model.User{
			Unit: &model.Unit{},
		}
		if err = rows.Scan(&ptr, &m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type); err != nil {
			c.Error = err
			return
		}

		request, ok := c.Requests[ptr]
		if !ok { // если еще не создавали то надо паниковать
			panic("c.Requests not exists")
		}
		request.Result.(*UserResult).User = m
	}

	c.Error = nil
}

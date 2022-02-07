package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *userResult) isRequestResult() {}
func (r *userInp) isRequestInput()     {}

type (
	userInp struct {
		UserID int
	}
	userResult struct {
		User *model.User
	}
)

func (d *Dataloader) User(userID int) (*model.User, error) {
	res := <-d.categories.User.addBaseRequest(
		&userInp{
			UserID: userID,
		},
		&userResult{
			User: &model.User{
				Unit: &model.Unit{},
			},
		},
	)
	if res == nil {
		return nil, d.categories.User.Error
	}
	return res.(*userResult).User, nil
}

func (c *parentCategory) user() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		userIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		userIDs = append(userIDs, query.Inp.(*userInp).UserID)
	}

	rows, err := c.Dataloader.db.Query(`
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
		request.Result.(*userResult).User = m
	}

	c.Error = nil
}

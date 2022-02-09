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
		new(userResult),
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
		SELECT ptr, 
		       coalesce(id, 0), 
		       coalesce(domain, ''), 
		       coalesce(name, ''), 
		       coalesce(type, 'USER') 
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, userid)
		LEFT JOIN units u ON u.id = inp.userid AND u.type = 'USER'
		`,
		pq.Array(ptrs),
		pq.Array(userIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("user:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr chanPtr
	)
	for rows.Next() {
		m := &model.User{Unit: new(model.Unit)}

		if err = rows.Scan(&ptr, &m.Unit.ID, &m.Unit.Domain, &m.Unit.Name, &m.Unit.Type); err != nil {
			//c.Dataloader.healer.Alert("user (scan rows):" + err.Error())
			c.Error = err
			return
		}
		if m.Unit.ID == 0 {
			m = nil
		}

		request := c.getRequest(ptr)
		request.Result.(*userResult).User = m
	}

	c.Error = nil
}

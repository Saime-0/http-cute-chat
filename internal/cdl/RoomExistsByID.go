package cdl

import (
	"fmt"
	"github.com/lib/pq"
)

func (r *roomExistsByIDResult) isRequestResult() {}
func (r *roomExistsByIDInp) isRequestInput()     {}

type (
	roomExistsByIDInp struct {
		RoomID int
	}
	roomExistsByIDResult struct {
		Exists bool
	}
)

func (d *Dataloader) RoomExistsByID(roomID int) (bool, error) {
	d.healer.Debug("Dataloader: новый запрос RoomExistsByID")
	res := <-d.categories.RoomExistsByID.addBaseRequest(
		&roomExistsByIDInp{
			RoomID: roomID,
		},
		new(roomExistsByIDResult),
	)
	if res == nil {
		return false, d.categories.RoomExistsByID.Error
	}
	return res.(*roomExistsByIDResult).Exists, nil
}

func (c *parentCategory) roomExistsByID() {
	var (
		inp = c.Requests

		ptrs    []chanPtr
		roomIDs []int
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		roomIDs = append(roomIDs, query.Inp.(*roomExistsByIDInp).RoomID)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, id is not null
		FROM unnest($1::varchar[], $2::bigint[]) inp(ptr, roomid)
		LEFT JOIN rooms u ON u.id = inp.roomid
		`,
		pq.Array(ptrs),
		pq.Array(roomIDs),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("roomExistsByID:" + err.Error())
		c.Error = err
		return
	}
	defer rows.Close()

	var ( // каждую итерацию будем менять значения
		ptr    chanPtr
		exists bool
	)
	for rows.Next() {

		if err = rows.Scan(&ptr, &exists); err != nil {
			//c.Dataloader.healer.Alert("roomExistsByID (scan rows):" + err.Error())
			c.Error = err
			return
		}

		request := c.getRequest(ptr)
		request.Result.(*roomExistsByIDResult).Exists = exists
	}

	c.Error = nil
}

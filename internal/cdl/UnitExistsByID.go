package cdl

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *unitExistsByIDResult) isRequestResult() {}
func (r *unitExistsByIDInp) isRequestInput()     {}

type (
	unitExistsByIDInp struct {
		UnitID   int
		UnitType model.UnitType
	}
	unitExistsByIDResult struct {
		Exists bool
	}
)

func (d *Dataloader) UnitExistsByID(unitID int, unitType model.UnitType) (bool, error) {
	res := <-d.categories.UnitExistsByID.addBaseRequest(
		&unitExistsByIDInp{
			UnitID:   unitID,
			UnitType: unitType,
		},
		new(unitExistsByIDResult),
	)
	if res == nil {
		return false, d.categories.UnitExistsByID.Error
	}
	return res.(*unitExistsByIDResult).Exists, nil
}

func (c *parentCategory) unitExistsByID() {
	var (
		inp = c.Requests

		ptrs      []chanPtr
		unitIDs   []int
		unitTypes []*model.UnitType
	)
	for _, query := range inp {
		ptrs = append(ptrs, fmt.Sprint(query.Ch))
		unitIDs = append(unitIDs, query.Inp.(*unitExistsByIDInp).UnitID)
		unitTypes = append(unitTypes, &query.Inp.(*unitExistsByIDInp).UnitType)
	}

	rows, err := c.Dataloader.db.Query(`
		SELECT ptr, id is not null
		FROM unnest($1::varchar[], $2::bigint[], $3::unit_type[]) inp(ptr, unitid, unittype)
		LEFT JOIN units u ON u.id = inp.unitid AND u.type = inp.unittype
		`,
		pq.Array(ptrs),
		pq.Array(unitIDs),
		pq.Array(unitTypes),
	)
	if err != nil {
		//c.Dataloader.healer.Alert("unitExistsByID:" + err.Error())
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
			//c.Dataloader.healer.Alert("unitExistsByID (scan rows):" + err.Error())
			c.Error = err
			return
		}

		request := c.getRequest(ptr)
		request.Result.(*unitExistsByIDResult).Exists = exists
	}

	c.Error = nil
}

package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/res"
)

type UnitsRepo struct {
	db *sql.DB
}

func NewUnitsRepo(db *sql.DB) *UnitsRepo {
	return &UnitsRepo{
		db: db,
	}
}

func (r *UnitsRepo) UnitExistsByID(unitId int, unitType res.UnitType) (exists bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units 
			WHERE id = $1 AND type = $2
			)`,
		unitId,
		unitType,
	).Scan(&exists)
	if err != nil {
		println("UnitExistsByID:", err.Error()) // debug
	}
	return
}

func (r *UnitsRepo) DomainIsFree(domain string) (free bool) {
	err := r.db.QueryRow(`
		SELECT 
		EXISTS (
			SELECT 1 
			FROM units
			WHERE domain = $1
		) 
		OR
		EXISTS (
		    SELECT 1
		    FROM registration_session
		    WHERE domain = $1
		)`,
		domain,
	).Scan(&free)
	if err != nil {
		println("DomainIsFree:", err.Error()) // debug
	}
	return !free
}

func (r *UnitsRepo) FindUnits(inp *model.FindUnits, params *model.Params) *model.Units {
	units := &model.Units{
		Units: []*model.Unit{},
	}
	if inp.NameFragment != nil {
		*inp.NameFragment = "%" + *inp.NameFragment + "%"
	}

	rows, err := r.db.Query(`
		SELECT id, domain, name , type 
		FROM units
		WHERE (
			    $1::BIGINT IS NULL 
			    OR id = $1
			)
			AND (
			    $2::VARCHAR IS NULL 
			    OR domain = $2 
			)
			AND (
			    $3::VARCHAR IS NULL 
			    OR name ILIKE $3
			)
			AND (
			    $4::char_type IS NULL 
			    OR type = $4
			)
		LIMIT $5
		OFFSET $6
		`,
		inp.ID,
		inp.Domain,
		inp.NameFragment,
		inp.UnitType,
		params.Limit,
		params.Offset,
	)
	defer rows.Close()
	if err != nil {
		println("FindUnits:", err.Error()) // debug
		return units
	}
	for rows.Next() {
		m := &model.Unit{}
		if err = rows.Scan(&m.ID, &m.Domain, &m.Name, &m.Type); err != nil {
			println("rows:Scan:", err.Error()) // debug
			return units
		}
		units.Units = append(units.Units, m)
	}
	return units
}

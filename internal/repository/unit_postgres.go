package repository

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

type UnitsRepo struct {
	db *sql.DB
}

func NewUnitsRepo(db *sql.DB) *UnitsRepo {
	return &UnitsRepo{
		db: db,
	}
}

func (r *UnitsRepo) GetDomainByID(unitId int) (domain string, err error) {
	return
}
func (r *UnitsRepo) GetIDByDomain(unitDomain string) (id int, err error) {
	panic("Not Implemented")
	return
}

func (r *UnitsRepo) UnitExistsByID(unitId int, unitType rules.UnitType) (exists bool) {
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
func (r *UnitsRepo) UnitExistsByDomain(unitDomain string, unitType rules.UnitType) (exists bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units
			WHERE domain = $1 AND type = $2
			)`,
		unitDomain,
		unitType,
	).Scan(&exists)
	if err != nil {
		println("UnitExistsByDomain:", err.Error()) // debug
	}
	return
}
func (r *UnitsRepo) DomainIsFree(domain string) (free bool) {
	err := r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units
			WHERE domain = $1
			)`,
		domain,
	).Scan(&free)
	if err != nil {
		println("DomainIsFree:", err.Error()) // debug
	}
	return !free
}

func (r *UnitsRepo) UnitByID(id int) (*model.Unit, error) {
	unit := &model.Unit{}
	err := r.db.QueryRow(`
		SELECT id, domain, name, type 
		FROM units
		WHERE id = $1`,
		id,
	).Scan(
		&unit.ID,
		&unit.Domain,
		&unit.Name,
		&unit.Type,
	)
	if err != nil {
		println("UnitByID:", err.Error()) // debug
	}
	return unit, err
}

func (r *UnitsRepo) FindUnits(inp *model.FindUnits, params *model.Params) *model.Units {
	units := &model.Units{}
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
		LIMIT $6
		OFFSET $7
		`,
		inp.ID,
		inp.Domain,
		inp.NameFragment,
		inp.UnitType,
	)
	if err != nil {
		println("FindUnits:", err.Error()) // debug
		return units
	}
	defer rows.Close()
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

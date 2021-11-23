package repository

import (
	"database/sql"
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
	return
}

func (r *UnitsRepo) UnitExistsByID(unitId int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units 
			WHERE id = $1
			)`,
		unitId,
	).Scan(&exists)
	return
}
func (r *UnitsRepo) UnitExistsByDomain(unitDomain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units
			WHERE domain = $1
			)`,
		unitDomain,
	).Scan(&exists)
	return
}

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

func (r *UnitsRepo) GetDomainByID(unit_id int) (domain string, err error) {
	return
}
func (r *UnitsRepo) GetIDByDomain(unit_domain string) (id int, err error) {
	return
}

func (r *UnitsRepo) UnitExistsByID(unit_id int) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units 
			WHERE id = $1
			)`,
		unit_id,
	).Scan(&exists)
	return
}
func (r *UnitsRepo) UnitExistsByDomain(unit_domain string) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units
			WHERE domain = $1
			)`,
		unit_domain,
	).Scan(&exists)
	return
}

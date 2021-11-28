package repository

import (
	"database/sql"
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
	return
}

func (r *UnitsRepo) UnitExistsByID(unitId int, unitType rules.UnitType) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units 
			WHERE id = $1 AND type = $2
			)`,
		unitId,
		unitType,
	).Scan(&exists)
	return
}
func (r *UnitsRepo) UnitExistsByDomain(unitDomain string, unitType rules.UnitType) (exists bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units
			WHERE domain = $1 AND type = $2
			)`,
		unitDomain,
		unitType,
	).Scan(&exists)
	return
}
func (r *UnitsRepo) DomainIsFree(domain string) (free bool) {
	r.db.QueryRow(
		`SELECT EXISTS(
			SELECT 1 
			FROM units
			WHERE domain = $1
			)`,
		domain,
	).Scan(&free)
	return
}

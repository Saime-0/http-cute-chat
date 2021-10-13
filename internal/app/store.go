package app

import "database/sql"

type Store struct {
	db *sql.DB
}

// New ...
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

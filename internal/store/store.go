package store

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/config"
)

func InitDB(cfg *config.Config2) (*sql.DB, error) {
	// connection string

	// open database
	db, err := sql.Open("postgres", cfg.PostgresConnection)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}

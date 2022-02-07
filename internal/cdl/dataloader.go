package cdl

import (
	"database/sql"
	"time"
)

type (
	ID      = string
	chanPtr = string
	Any     = interface{}

	RequestsCount uint8
	//categories    map[CategoryName]*parentCategory
)

type Dataloader struct {
	wait             time.Duration
	capactiyRequests RequestsCount
	categories       *Categories
	db               *sql.DB
}

func NewDataloader(wait time.Duration, maxBatch RequestsCount, db *sql.DB) *Dataloader {
	d := &Dataloader{
		wait:             wait,
		capactiyRequests: maxBatch,
		db:               db,
	}
	d.ConfigureDataloader()
	return d
}

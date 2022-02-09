package cdl

import (
	"database/sql"
	"github.com/saime-0/http-cute-chat/internal/healer"
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
	healer           *healer.Healer
}

func NewDataloader(wait time.Duration, maxBatch RequestsCount, db *sql.DB, hlr *healer.Healer) *Dataloader {
	d := &Dataloader{
		wait:             wait,
		capactiyRequests: maxBatch,
		db:               db,
		healer:           hlr,
	}
	d.ConfigureDataloader()
	return d
}

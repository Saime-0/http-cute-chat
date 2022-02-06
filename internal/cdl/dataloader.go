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
	//Categories    map[CategoryName]*ParentCategory
)

type Dataloader struct {
	Wait             time.Duration
	CapactiyRequests RequestsCount
	Categories       *Categories
	DB               *sql.DB
}

func NewDataloader(wait time.Duration, maxBatch RequestsCount, db *sql.DB) *Dataloader {
	d := &Dataloader{
		Wait:             wait,
		CapactiyRequests: maxBatch,
		DB:               db,
	}
	d.ConfigureDataloader()
	return d
}

func (d *Dataloader) ConfigureDataloader() {
	d.Categories = &Categories{
		Rooms:            d.NewRoomsCategory(),
		UserIsChatMember: d.NewUserIsChatMemberCategory(),
	}
}

//func (d *Dataloader) AddCategory(v CategoryCell) {
//	d.Categories[v.Name] = &ParentCategory{
//		Dataloader:             d,
//		RemainingRequestsCount: d.CapactiyRequests,
//		//State:                  SLEEP,
//		//StateCh:                make(chan CategoryState),
//		LoadFn: v.LoadFn,
//	}
//}

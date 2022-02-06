package cdl

import (
	"fmt"
	"sync"
	"time"
)

type Categories struct {
	Rooms            *ParentCategory
	UserIsChatMember *ParentCategory
}

type ParentCategory struct {
	Dataloader             *Dataloader
	RemainingRequestsCount RequestsCount
	Timer                  *time.Timer

	Error error

	LoadFn func()
	//PrepareForNextLaunch func()

	Requests map[chanPtr]*BaseRequest
	sync.Mutex
}

//func (n Null) IsRequestResult()  {}
//type Null types.Nil

func (d *Dataloader) NewParentCategory() *ParentCategory {
	c := &ParentCategory{
		Dataloader:             d,
		RemainingRequestsCount: d.CapactiyRequests,
		Requests:               map[chanPtr]*BaseRequest{},
		LoadFn: func() {
			panic("not implemented")
		},
	}
	return c
}

func (c *ParentCategory) RunLoadFunc() {
	c.LoadFn()
	var result RequestResult
	c.Lock() // Lock
	for _, request := range c.Requests {
		if c.Error != nil {
			result = nil
		} else {
			result = request.Result
		}
		select {
		case request.Ch <- result:
		default:
		}
	}
	// PrepareForNextLaunch
	for ptr := range c.Requests {
		delete(c.Requests, ptr)
	}
	c.RemainingRequestsCount = c.Dataloader.CapactiyRequests
	c.Timer = nil
	c.Unlock() // Unlock
}

// deprecated
func (c *ParentCategory) PrepareForNextLaunchV2() {
	c.Lock()
	for ptr := range c.Requests {
		delete(c.Requests, ptr)
	}
	c.Unlock()
}

func (c *ParentCategory) OnAddRequest() {
	fmt.Println("реквест добавлен") // debug
	fmt.Printf("таймер на момент OnAddRequest %#v\n", c.Timer)
	c.RemainingRequestsCount -= 1
	if c.RemainingRequestsCount <= 0 {
		if c.Timer != nil {
			c.Timer.Stop()
			fmt.Println("таймер остановлен тк нужное количество реквестов набралось") // debug
		}
		c.RunLoadFunc()
		fmt.Println("функция выполнилась по набору максимального количества реквестов") // debug
		return
	}

	if c.Timer == nil {
		fmt.Println("таймер запущен") // debug
		c.Timer = time.AfterFunc(c.Dataloader.Wait, func() {
			c.RunLoadFunc()
			fmt.Println("функция выполнилась по таймеру") // debug
		})
	}
}

//func (c *ParentCategory) RunLoadFnByTrigger() {
//	c.LoadFn()
//	c.RemainingRequestsCount = c.Dataloader.CapactiyRequests
//	c.Timer = nil
//	c.PrepareForNextLaunch()
//}

//func (b BaseResultChan) IsResultChan() {}

type BaseResultChan chan RequestResult

func (b BaseRequest) IsRequest() {}

type BaseRequest struct {
	Ch     BaseResultChan
	Inp    RequestInput
	Result RequestResult
}

type (
	Request interface {
		IsRequest()
		//GetChan() BaseResultChan
		//GetResult() RequestResult
	}
	RequestResult interface{ IsRequestResult() }
	RequestInput  interface{ IsRequestInput() }
	//ResultChan    interface{ IsResultChan() }
)

//type RequestsMapValue interface {
//	is
//}

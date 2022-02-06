package cdl

import (
	"fmt"
	"time"
)

type Categories struct {
	Rooms            *RoomsCategory
	UserIsChatMember *UserIsChatMemberCategory
}

type ParentCategory struct {
	Dataloader             *Dataloader
	RemainingRequestsCount RequestsCount
	Timer                  *time.Timer

	Error error

	LoadFn               func()
	PrepareForNextLaunch func()
}

//type RequestResult interface {
//	IsRequestResult()
//}
//
//type Category interface {
//	IsCategory()
//	//AddRequest() chan *RequestResult
//}

func (c *ParentCategory) OnAddRequest() {
	fmt.Println("реквест добавлен") // debug

	c.RemainingRequestsCount -= 1
	if c.RemainingRequestsCount <= 0 {
		if c.Timer != nil {
			c.Timer.Stop()
			fmt.Println("таймер остановлен тк нужное количество реквестов набралось") // debug
		}
		c.RunLoadFnByTrigger()
		fmt.Println("функция выполнилась по набору максимального количества реквестов") // debug
		return
	}

	if c.Timer == nil {
		c.Timer = time.AfterFunc(c.Dataloader.Wait, func() {
			c.RunLoadFnByTrigger()
			fmt.Println("функция выполнилась по таймеру") // debug
		})
	}
}

func (c *ParentCategory) RunLoadFnByTrigger() {
	c.LoadFn()
	c.RemainingRequestsCount = c.Dataloader.CapactiyRequests
	c.Timer = nil
	c.PrepareForNextLaunch()
}

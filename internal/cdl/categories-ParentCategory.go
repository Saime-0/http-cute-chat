package cdl

import (
	"fmt"
	"sync"
	"time"
)

type BaseResultChan chan RequestResult

func (b BaseRequest) IsRequest() {}

type BaseRequest struct {
	Ch     BaseResultChan
	Inp    RequestInput
	Result RequestResult
}

type (
	Request       interface{ IsRequest() }
	RequestResult interface{ IsRequestResult() }
	RequestInput  interface{ IsRequestInput() }
)

type ParentCategory struct {
	Dataloader             *Dataloader
	RemainingRequestsCount RequestsCount
	Timer                  *time.Timer
	Error                  error
	LoadFn                 func()
	Requests               map[chanPtr]*BaseRequest
	sync.Mutex
}

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
	c.Lock() // Lock
	c.LoadFn()
	var result RequestResult
	for _, request := range c.Requests {
		// почему проверка ошибки, а не отправка request.Result?
		// все просто, если в load функции или раньше, будет
		// подготовлен request то смысла здесь провреки на nil нет
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

func (c *ParentCategory) AddBaseRequest(inp RequestInput, result RequestResult) BaseResultChan {

	newClient := make(BaseResultChan)
	c.Lock()
	c.Requests[fmt.Sprint(newClient)] = &BaseRequest{
		Ch:     newClient,
		Inp:    inp,
		Result: result,
	}
	c.Unlock()

	go c.OnAddRequest()
	return newClient

}

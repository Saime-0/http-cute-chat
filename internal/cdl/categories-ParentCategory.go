package cdl

import (
	"fmt"
	"sync"
	"time"
)

type baseResultChan chan requestResult

func (b baseRequest) isRequest() {}

type baseRequest struct {
	Ch     baseResultChan
	Inp    requestInput
	Result requestResult
}

type (
	requestResult interface{ isRequestResult() }
	requestInput  interface{ isRequestInput() }
)

type parentCategory struct {
	Dataloader             *Dataloader
	RemainingRequestsCount RequestsCount
	Timer                  *time.Timer
	Error                  error
	LoadFn                 func()
	Requests               map[chanPtr]*baseRequest
	sync.Mutex
}

func (d *Dataloader) newParentCategory() *parentCategory {
	c := &parentCategory{
		Dataloader:             d,
		RemainingRequestsCount: d.capactiyRequests,
		Requests:               map[chanPtr]*baseRequest{},
		LoadFn: func() {
			panic("not implemented")
		},
	}
	return c
}

func (c *parentCategory) runLoadFunc() {
	c.Lock() // Lock
	c.LoadFn()
	var result requestResult
	for _, request := range c.Requests {
		// почему проверка ошибки, а не отправка request.Result?
		// все просто, если в load функции или раньше, будет
		// подготовлен request то смысла здесь провреки на nil нет
		// upd: nil, 0, "" так же могу отправлять и это будет означать что объект не найден
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
	c.RemainingRequestsCount = c.Dataloader.capactiyRequests
	c.Timer = nil
	c.Unlock() // Unlock
}

func (c *parentCategory) onAddRequest() {
	// реквест уже добавлен

	c.RemainingRequestsCount -= 1
	if c.RemainingRequestsCount <= 0 {
		if c.Timer != nil {
			// таймер остановлен тк нужное количество реквестов набралось
			c.Timer.Stop()
		}
		// функция выполнилнится по набору максимального количества реквестов
		c.runLoadFunc()
		return
	}

	if c.Timer == nil {
		// запуск таймера
		c.Timer = time.AfterFunc(c.Dataloader.wait, func() {
			c.runLoadFunc()
		})
	}
}

func (c *parentCategory) addBaseRequest(inp requestInput, result requestResult) baseResultChan {

	newClient := make(baseResultChan)
	c.Lock()
	c.Requests[fmt.Sprint(newClient)] = &baseRequest{
		Ch:     newClient,
		Inp:    inp,
		Result: result,
	}
	c.Unlock()

	go c.onAddRequest()
	return newClient

}

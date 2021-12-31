package mutexmaps

import (
	"sync"
)

type subGroup struct {
	Root *Subscriber
	mu   *sync.Mutex
}

func newGroup() *subGroup {
	return &subGroup{
		Root: &Subscriber{},
		mu:   new(sync.Mutex),
	}
}

type Subscriber struct {
	Ch    interface{}
	group *subGroup
	next  *Subscriber
	prev  *Subscriber
}

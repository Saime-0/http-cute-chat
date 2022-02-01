package scheduler

import (
	"sync"
	"time"
)

type SheduleStatus int8

const (
	RUNNING SheduleStatus = iota
	IDLE
)

type Scheduler struct {
	mu     *sync.Mutex
	root   *Task
	Status SheduleStatus
	Ch     chan SheduleStatus
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		mu: new(sync.Mutex),
		root: &Task{
			next: emptyTask(),
		},
		Status: IDLE,
		Ch:     make(chan SheduleStatus),
	}
	go s.proccess()
	<-s.Ch
	return s
}

func (s *Scheduler) Resume() {
	s.Status = RUNNING
	select {
	case s.Ch <- RUNNING:
	default: // "deadlock" fix
	}

}
func (s *Scheduler) Pause() {
	s.Status = IDLE
	select {
	case s.Ch <- IDLE:
	default: // "deadlock" fix
	}
}

func (s *Scheduler) proccess() {
	s.Ch <- RUNNING
	for {
		for s.Status == RUNNING {

			if *s.root.next == nil {
				s.Status = IDLE
				break
			} else if (*s.root.next).executeAt <= time.Now().Unix() {
				s.runNextTask()
				continue
			}
			time.Sleep(time.Second)
		}

		<-s.Ch
	}
}

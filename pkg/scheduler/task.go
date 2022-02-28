package scheduler

import (
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"log"
	"time"
)

type execTaskFn func()

type Task struct {
	sch       *Scheduler
	executeAt int64
	taskFunc  execTaskFn
	next      **Task
}

func emptyTask() **Task {
	var x *Task
	return &x
}

func (s *Scheduler) RunTask(task **Task) error {
	// принудительный запуск задачи
	s.mu.Lock()
	defer s.mu.Unlock()
	if task == nil || *task == nil {
		return cerrors.New("task not found")
	}
	go (*task).taskFunc()
	*task = *(*task).next
	return nil
}

func (s *Scheduler) runNextTask() {
	s.mu.Lock()
	task := s.root.next
	log.Println("runNextTask: запускаю задачу..") // debug
	go (*task).taskFunc()
	*task = *(*task).next
	s.mu.Unlock()
}

func (s *Scheduler) AddTask(taskFn execTaskFn, executeAt int64) (newTask *Task, err error) {
	if executeAt <= time.Now().Unix() {
		return nil, cerrors.New("not valid executeAt")
	}
	newTask = &Task{
		sch:       s,
		executeAt: executeAt,
		taskFunc:  taskFn,
		next:      emptyTask(),
	}
	task := s.root
	s.mu.Lock()

	for *task.next != nil {
		if task.executeAt < newTask.executeAt && (*task.next).executeAt >= newTask.executeAt {
			newTask.next = task.next
			break
		}
		task = *task.next
	}
	task.next = &newTask
	s.mu.Unlock()
	if task == s.root {
		s.Resume()
	}
	return
}

func (s *Scheduler) DropTask(task **Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task == nil || *task == nil {
		return cerrors.New("task not found")
	}
	*task = *(*task).next
	return nil
}

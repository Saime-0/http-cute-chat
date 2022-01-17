package scheduler

import (
	"errors"
	"fmt"
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if task == nil || *task == nil {
		return errors.New("task not found")
	}
	go (*task).taskFunc()
	*task = *(*task).next
	println("RunTask: задача запущена принудительно") // debug
	return nil
}

func (s *Scheduler) runNextTask() {
	s.mu.Lock()
	task := s.root.next
	go (*task).taskFunc()
	*task = *(*task).next
	s.mu.Unlock()
	println("runNextTask: задача запущена") // debug
}

func (s *Scheduler) AddTask(taskFn execTaskFn, executeAt int64) (newTask *Task, err error) {
	if executeAt <= time.Now().Unix() {
		return nil, errors.New("not valid executeAt")
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
		fmt.Printf("Task at %d, next: %#v |", task.executeAt, task.next) // debug
		if task.executeAt < newTask.executeAt && (*task.next).executeAt >= newTask.executeAt {
			newTask.next = task.next
			break
		}
		task = *task.next
	}
	println()
	task.next = &newTask
	s.mu.Unlock()

	if task == s.root {
		s.Resume()
	}
	println("AddTask: добавлена новая задача") // debug
	return
}

func (s *Scheduler) DropTask(task **Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task == nil || *task == nil {
		return errors.New("task not found")
	}
	*task = *(*task).next
	println("DropTask: задача удалена") // debug
	return nil
}

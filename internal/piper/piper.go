package piper

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Pipeline struct {
	RootNode *Node
	repos    *repository.Repositories
}

func NewPipeline(repos *repository.Repositories) *Pipeline {
	return &Pipeline{
		RootNode: &Node{},
		repos:    repos,
	}
}

type Node struct {
	prev     *Node
	next     *Node
	repos    *repository.Repositories
	Err      *model.AdvancedError
	timeline *TimeLine
}

func (p *Pipeline) CreateNode(processName ...interface{}) *Node {
	n := &Node{
		prev:     p.RootNode,
		next:     p.RootNode.next,
		repos:    p.repos,
		timeline: Start(processName),
	}
	if p.RootNode.next != nil {
		p.RootNode.next.prev = n
	}
	p.RootNode.next = n
	return n
}

func (n *Node) Kill() *model.AdvancedError {
	n.prev.next = n.next
	if n.next != nil {
		n.next.prev = n.prev
	}
	n.timeline.Measure()
	return n.Err
}

func (n *Node) printProcessline(timerID int) {

}

// todo у ноды должно быть поле с таймлайном, у поинтов таймлайна поля ид и время начала
// и методы: завершить поинт и вывести таймлайн

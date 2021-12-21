package piper

import "github.com/saime-0/http-cute-chat/graph/model"

type Pipeline struct {
	RootNode *Node
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		RootNode: &Node{},
	}
}

type Node struct {
	prev *Node
	next *Node
	Err  *model.AdvancedError
}

func (p *Pipeline) CreateNode() *Node {
	n := &Node{
		prev: p.RootNode,
		next: p.RootNode.next,
	}
	if p.RootNode.next != nil {
		p.RootNode.next.prev = n
	}
	p.RootNode.next = n
	return n
}

func (n *Node) Kill() *model.AdvancedError {
	n.prev.next, n.next.prev = n.next, n.prev
	return n.Err
}

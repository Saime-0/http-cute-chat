package piper

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Rows []interface{}

type Request struct {
	Timestamp time.Time
	ID        string
	Status    int
	Method    string
	Path      string
	Duration  string
	Body      *Rows
}

type Method struct {
	Method   string
	startAt  time.Time
	Duration string
	Body     *Rows
}

type Node struct {
	repos  *repository.Repositories
	Healer *healer.Healer
	Err    *model.AdvancedError

	ID            *string
	RootContainer interface{}
	scope         *Rows
	ScopeMethod   *Method
}

func (p *Pipeline) CreateNode(id string) (*Node, *Request) {
	scope := &Rows{}

	request := &Request{
		Timestamp: time.Now(),
		ID:        kit.RandomSecret(6),
		Body:      scope,
	}

	n := &Node{
		repos:  p.repos,
		Healer: p.healer,

		ID: &id,
		RootContainer: bson.M{
			"Request": request,
		},
		scope: scope,
	}
	p.Nodes[id] = n
	return n, request
}

func (n *Node) Delete() {
	panic("deprecated")
}

func (p *Pipeline) DeleteNode(id string) {
	delete(p.Nodes, id)
}

func (n *Node) Execute() {
	n.Healer.Log(n.RootContainer)
}

func (n *Node) SwitchMethod(name string) {
	meth := &Method{
		Method:  name,
		Body:    &Rows{},
		startAt: time.Now(),
	}
	*n.scope = append(*n.scope, meth)
	n.scope = meth.Body

	n.ScopeMethod = meth
}
func (n *Node) MethodTiming() {
	if n.ScopeMethod != nil {
		n.ScopeMethod.Duration = time.Since(n.ScopeMethod.startAt).String()
	}
}
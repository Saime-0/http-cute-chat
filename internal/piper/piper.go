package piper

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/cdl"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/res"
)

type Pipeline struct {
	Nodes map[string]*Node

	repos      *repository.Repositories
	healer     *healer.Healer
	dataloader *cdl.Dataloader
}

func NewPipeline(repos *repository.Repositories, healer *healer.Healer, dataloader *cdl.Dataloader) *Pipeline {
	return &Pipeline{
		Nodes:      map[string]*Node{},
		repos:      repos,
		healer:     healer,
		dataloader: dataloader,
	}
}

func (p *Pipeline) NodeFromContext(ctx context.Context) *Node {
	return ctx.Value(res.CtxNode).(*Node)
}

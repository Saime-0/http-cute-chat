package piper

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/cdl"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/res"
)

type Pipeline struct {
	Nodes map[string]*Node

	repos      *repository.Repositories
	healer     *healer.Healer
	dataloader *cdl.Dataloader

	cfg *config.Config2
}

func NewPipeline(cfg *config.Config2, repos *repository.Repositories, healer *healer.Healer, dataloader *cdl.Dataloader) *Pipeline {
	return &Pipeline{
		Nodes:      map[string]*Node{},
		repos:      repos,
		healer:     healer,
		dataloader: dataloader,
		cfg:        cfg,
	}
}

func (p *Pipeline) NodeFromContext(ctx context.Context) *Node {
	return ctx.Value(res.CtxNode).(*Node)
}

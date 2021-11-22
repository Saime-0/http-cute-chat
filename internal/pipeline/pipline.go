package pipeline

import (
	"context"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Pipeline struct {
	Ctx   context.Context
	Repos *repository.Repositories
	Err   *model.AdvancedError
}

func NewPipeline(ctx context.Context, repos *repository.Repositories) *Pipeline {
	return &Pipeline{
		Ctx:   ctx,
		Repos: repos,
		Err:   nil,
	}
}

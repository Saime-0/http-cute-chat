package piping

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Pipeline struct {
	Ctx      context.Context
	repos    *repository.Repositories
	overcook overcooker
	Err      *model.AdvancedError
	Can      Can
}

func NewPipeline(ctx context.Context, repos *repository.Repositories) *Pipeline {
	// надо бы какую нибудь штуку придумуть для возврата значений из  обработчиков
	p := &Pipeline{
		Ctx:      ctx,
		repos:    repos,
		overcook: overcooker{},
		Err:      nil,
	}
	p.Can.pl = p
	return p
}

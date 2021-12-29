package piping

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/repository"
)

type Pipeline struct {
	repos *repository.Repositories
	Err   *model.AdvancedError
	Can   Can
}

func NewPipeline(repos *repository.Repositories) *Pipeline {
	// надо бы какую нибудь штуку придумуть для возврата значений из  обработчиков
	p := &Pipeline{
		repos: repos,
		Err:   nil,
	}
	p.Can.pl = p
	return p
}

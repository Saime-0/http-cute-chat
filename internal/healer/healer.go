package healer

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/service"
	"github.com/saime-0/http-cute-chat/pkg/fsm"
)

type Healer struct {
	stateMachine *fsm.Machine
	services     *service.Services
	cfg          *config.Config
}

func NewHealer(s *service.Services, cfg *config.Config) (*Healer, error) {
	machine, err := fsm.NewMachine()
	if err != nil {
		return nil, errors.Wrap(err, res.FailedToCreateHealer)
	}
	h := &Healer{
		stateMachine: machine,
		services:     s,
		cfg:          cfg,
	}

	err = h.prepareHealer()
	if err != nil {
		return nil, errors.Wrap(err, res.FailedToPrepareHealer)
	}

	return h, nil
}

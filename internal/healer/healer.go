package healer

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/cache"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/pkg/fsm"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
)

type Healer struct {
	stateMachine *fsm.Machine
	cfg          *config.Config
	sched        *scheduler.Scheduler
	cache        *cache.Cache
	logger       *clog.Clog
}

func NewHealer(cfg *config.Config, sched *scheduler.Scheduler, cache *cache.Cache, logger *clog.Clog) (*Healer, error) {
	machine, err := fsm.NewMachine()
	if err != nil {
		return nil, errors.Wrap(err, res.FailedToCreateHealer)
	}
	h := &Healer{
		stateMachine: machine,
		cfg:          cfg,
		sched:        sched,
		cache:        cache,
		logger:       logger,
	}

	err = h.prepareHealer()
	if err != nil {
		return nil, errors.Wrap(err, res.FailedToPrepareHealer)
	}

	return h, nil
}

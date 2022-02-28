package healer

import (
	"github.com/saime-0/http-cute-chat/internal/cache"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/pkg/fsm"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ clog.Logger = (*Healer)(nil)

type Healer struct {
	stateMachine *fsm.Machine
	cfg          *config.Config2
	sched        *scheduler.Scheduler
	cache        *cache.Cache

	// logging
	db     *mongo.Database
	Level  clog.LogLevel
	Output clog.Output
	client *mongo.Client
}

func NewHealer(cfg *config.Config2, sched *scheduler.Scheduler, cache *cache.Cache) (*Healer, error) {
	machine, err := fsm.NewMachine()
	if err != nil {
		return nil, cerrors.Wrap(err, res.FailedToCreateHealer)
	}
	h := &Healer{
		stateMachine: machine,
		cfg:          cfg,
		sched:        sched,
		cache:        cache,
	}

	if err := h.prepareHealer(); err != nil {
		return nil, cerrors.Wrap(err, res.FailedToPrepareHealer)
	}
	if err := h.PrepareLogging(cfg); err != nil {
		return nil, cerrors.Wrap(err, "не удалось настроить логирование")
	}
	return h, nil
}

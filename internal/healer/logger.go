package healer

import (
	"context"
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	CouldNotConnectToDb = "could not connect to db"
	ConnectionTimeout   = time.Millisecond * 1500
)

func (h *Healer) PrepareLogging(cfg *config.Config) (err error) {

	lvl, err := clog.GetLogLevel(cfg.Logger.Level)
	if err != nil {
		return errors.Wrap(err, "не удалось определить уровень логирования")
	}

	h.Level = lvl
	h.Output = cfg.Logger.Output

	if h.Output < clog.Multiple {
		return nil
	} // don't creating db connection

	clientOptions := options.Client().ApplyURI("mongodb+srv://" +
		cfg.Logger.MongoDBUser +
		":" + cfg.Logger.MongoDBPassword +
		"@" + cfg.Logger.MongoDBCluster,
	)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return errors.Wrap(err, CouldNotConnectToDb)
	}

	db := client.Database(cfg.Logger.DBName)

	if err != nil {
		return err
	}

	h.db = db
	h.client = client

	return nil
}

func (h *Healer) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()
	return h.client.Ping(ctx, readpref.Primary())
}

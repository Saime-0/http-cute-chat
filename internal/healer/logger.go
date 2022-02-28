package healer

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
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

func (h *Healer) PrepareLogging(cfg *config.Config2) (err error) {

	h.Level = clog.LogLevel(*cfg.Logging.LoggingLevel)
	h.Output = clog.Output(*cfg.Logging.LoggingOutput)
	if h.Output < clog.Multiple {
		return nil
	} // don't creating db connection

	clientOptions := options.Client().ApplyURI(cfg.MongoDBUri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return cerrors.Wrap(err, CouldNotConnectToDb)
	}

	db := client.Database(*cfg.Logging.LoggingDBName)

	if err != nil {
		return err
	}

	h.db = db
	h.client = client
	if err := h.PingDB(); err != nil {
		return cerrors.Wrap(err, "неудачное подключение mongodb")
	}

	return nil
}

func (h *Healer) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()
	return h.client.Ping(ctx, readpref.Primary())
}

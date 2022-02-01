package clog

import (
	"context"
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Clog struct {
	db     *mongo.Database
	Level  LogLevel
	Output Output
	client *mongo.Client
}

type Output uint8

const (
	_ Output = iota
	Console
	Multiple
	MongoDB
)

func NewClog(cfg *config.Config, output Output) (*Clog, error) {

	lvl := GetLogLevel(cfg.Logger.Level)
	if lvl == -1 {
		return nil, errors.New("the required level does not exist")
	}

	clientOptions := options.Client().ApplyURI("mongodb+srv://" +
		cfg.Logger.MongoDBUser +
		":" + cfg.Logger.MongoDBPassword +
		"@" + cfg.Logger.MongoDBCluster,
	)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, CouldNotConnectToDb)
	}

	db := client.Database(cfg.Logger.DBName)

	if err != nil {
		return nil, err
	}

	c := &Clog{
		db:     db,
		Level:  lvl,
		Output: output,
		client: client,
	}

	return c, nil
}

const ConnectionTimeout = time.Millisecond * 1500
const CouldNotConnectToDb = "could not connect to db"

func (c *Clog) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()
	return c.client.Ping(ctx, readpref.Primary())
}

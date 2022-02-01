package clog

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Clog struct {
	db     *mongo.Database
	Level  LogLevel
	Output Output
	// for reconnecting
	clientOptions *options.ClientOptions
	dbName        string
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
	db, err := connectToDB(
		cfg.Logger.DBName,
		clientOptions,
	)
	if err != nil {
		return nil, err
	}

	c := &Clog{
		db:            db,
		Level:         lvl,
		Output:        output,
		clientOptions: clientOptions,
	}

	return c, nil
}

func (c Clog) ReconnectToDB() (err error) {
	db, err := connectToDB(c.dbName, c.clientOptions)
	if err != nil {
		return errors.Wrap(err, "could not reconnect to db")
	}
	c.db = db
	return nil
}

func connectToDB(dbName string, clientOptions *options.ClientOptions) (*mongo.Database, error) {
	client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to db")
	}
	return client.Database(dbName), nil
}

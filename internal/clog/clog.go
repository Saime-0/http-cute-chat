package clog

import (
	"errors"
	"github.com/saime-0/http-cute-chat/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Clog struct {
	db    *mongo.Database
	Level LogLevel
}

func NewClog(cfg *config.Config) (*Clog, error) {

	lvl := GetLogLevel(cfg.Logger.Level)
	if lvl == -1 {
		return nil, errors.New("the required level does not exist")
	}

	clientOptions := options.Client().ApplyURI("mongodb+srv://" +
		cfg.Logger.MongoDBUser +
		":" + cfg.Logger.MongoDBPassword +
		"@" + cfg.Logger.MongoDBCluster,
	)

	client, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.Logger.DBName)

	c := &Clog{
		db:    db,
		Level: lvl,
	}

	return c, nil
}

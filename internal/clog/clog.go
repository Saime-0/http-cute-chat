package clog

import (
	"go.mongodb.org/mongo-driver/mongo"
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

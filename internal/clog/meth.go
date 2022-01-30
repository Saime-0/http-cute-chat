package clog

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c Clog) log(lvl LogLevel, document interface{}) (err error) {
	doc := bson.M{
		"time":    time.Now().Unix(),
		"type":    lvl.String(),
		"payload": document,
	}
	_, err = c.db.Collection("logs").InsertOne(nil, doc)

	return err
}

func (c Clog) Emergency(document interface{}) (err error) {
	if c.Level >= Emergency {
		err = c.log(Emergency, document)
	}
	return
}
func (c Clog) Alert(document interface{}) (err error) {
	if c.Level >= Emergency {
		err = c.log(Alert, document)
	}
	return
}
func (c Clog) Critical(document interface{}) (err error) {
	if c.Level >= Critical {
		err = c.log(Critical, document)
	}
	return
}
func (c Clog) Error(document interface{}) (err error) {
	if c.Level >= Error {
		err = c.log(Error, document)
	}
	return
}
func (c Clog) Warning(document interface{}) (err error) {
	if c.Level >= Warning {
		err = c.log(Warning, document)
	}
	return
}
func (c Clog) Notice(document interface{}) (err error) {
	if c.Level >= Notice {
		err = c.log(Notice, document)
	}
	return
}
func (c Clog) Info(document interface{}) (err error) {
	if c.Level >= Info {
		err = c.log(Info, document)
	}
	return
}
func (c Clog) Debug(document interface{}) (err error) {
	if c.Level >= Debug {
		err = c.log(Debug, document)
	}
	return
}

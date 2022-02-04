package clog

type Logger interface {
	Emergency(document interface{})
	Alert(document interface{})
	Critical(document interface{})
	Error(document interface{})
	Notice(document interface{})
	Info(document interface{})
	Debug(document interface{})
}

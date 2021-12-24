package tlog

import (
	"fmt"
	"time"
)

const (
	FINE = true
	TIME = true
)

type Logger struct {
	processName string
	start       time.Time
}

func Start(pname ...interface{}) *Logger {
	return &Logger{
		processName: fmt.Sprint(pname),
		start:       time.Now(),
	}
}
func (l *Logger) Time() {
	printTime(l, "")
}
func (l *Logger) TimeWithStatus(status ...interface{}) {
	printTime(l, "\t", status)
}
func (l *Logger) Fine() {
	printFine(l, "")
}
func (l *Logger) FineWithReason(reason ...interface{}) {
	printFine(l, reason)
}
func printFine(l *Logger, reason ...interface{}) {
	if FINE {
		fmt.Println("TLOG:", l.processName, "| finally", reason, "| duration", time.Since(l.start))
	}
}
func printTime(l *Logger, status ...interface{}) {
	if TIME {
		fmt.Println(l.processName, "-", time.Since(l.start), status)
	}
}

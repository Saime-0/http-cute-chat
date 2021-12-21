package tlog

import (
	"fmt"
	"time"
)

const (
	FINE = false
	TIME = true
)

type Logger struct {
	processName string
	start       time.Time
}

func Start(pname string) *Logger {
	return &Logger{
		processName: pname,
		start:       time.Now(),
	}
}
func (l *Logger) Time() {
	printTime(l, "")
}
func (l *Logger) TimeWithStatus(status string) {
	printTime(l, "\t"+status)
}
func (l *Logger) Fine() {
	printFine(l, "")
}
func (l *Logger) FineWithReason(reason string) {
	printFine(l, reason)
}
func printFine(l *Logger, reason string) {
	if FINE {
		fmt.Println("TLOG:", l.processName, "| finally", reason, "| duration", time.Since(l.start))
	}
}
func printTime(l *Logger, status string) {
	if TIME {
		fmt.Println(l.processName, "-", time.Since(l.start), status)
	}
}

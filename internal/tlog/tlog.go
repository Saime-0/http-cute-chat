package tlog

import (
	"fmt"
	"time"
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
func (l *Logger) Fine() {
	fmt.Println("P:", l.processName, "finally duration", time.Since(l.start))
}
func (l *Logger) FineWithReason(reason string) {
	fmt.Println("P:", l.processName, "finally(", reason, ") duration", time.Since(l.start))
}

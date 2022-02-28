package clog

import "github.com/saime-0/http-cute-chat/internal/cerrors"

type LogLevel uint8

const (
	Emergency LogLevel = iota // system is unusable
	Alert                     // action must be taken immediately
	Critical                  // critical conditions
	Error                     // error conditions
	Warning                   // warning conditions
	Notice                    // normal but significant condition
	Info                      // informational messages
	Debug                     // debug-level messages
)

var LogLevelNotExists = cerrors.New("the required level does not exist")
var lvlNames = []string{
	"emergency",
	"alert",
	"critical",
	"error",
	"warning",
	"notice",
	"info",
	"debug",
}

func (lvl LogLevel) String() string {
	return lvlNames[lvl]
}

func GetLogLevel(str string) (lvl LogLevel, err error) {
	for lvl, name := range lvlNames {
		if name == str {
			return LogLevel(lvl), nil
		}
	}
	return 0, LogLevelNotExists
}

func Exists(lvl LogLevel) bool {
	if lvl < 1 || int(lvl) >= len(lvlNames) {
		return false
	}
	return true
}

package clog

type LogLevel int8

const (
	_         LogLevel = iota - 1
	Emergency          // system is unusable
	Alert              // action must be taken immediately
	Critical           // critical conditions
	Error              // error conditions
	Warning            // warning conditions
	Notice             // normal but significant condition
	Info               // informational messages
	Debug              // debug-level messages
)

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

func GetLogLevel(str string) LogLevel {
	for lvl, name := range lvlNames {
		if name == str {
			return LogLevel(lvl)
		}
	}
	return -1
}

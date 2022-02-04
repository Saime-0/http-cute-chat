package piper

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/internal/clog"
)

type LoggingRow struct {
	Level string
	Body  interface{}
}

func NewLoggingRow(lvl string, body interface{}) *LoggingRow {
	return &LoggingRow{
		Level: lvl,
		Body:  body,
	}
}

func (n *Node) addLoggingRowToScope(level clog.LogLevel, document interface{}) {

	if n.scope != nil {
		*n.scope = append(*n.scope, NewLoggingRow(
			level.String(),
			fmt.Sprint(document),
		))
	}
}

func (n *Node) Emergency(document interface{}) {
	n.addLoggingRowToScope(clog.Emergency, document)
}
func (n *Node) Alert(document interface{}) {
	n.addLoggingRowToScope(clog.Alert, document)
}
func (n *Node) Critical(document interface{}) {
	n.addLoggingRowToScope(clog.Critical, document)
}
func (n *Node) Error(document interface{}) {
	n.addLoggingRowToScope(clog.Error, document)
}
func (n *Node) Notice(document interface{}) {
	n.addLoggingRowToScope(clog.Notice, document)
}
func (n *Node) Info(document interface{}) {
	n.addLoggingRowToScope(clog.Info, document)
}
func (n *Node) Debug(document interface{}) {
	n.addLoggingRowToScope(clog.Debug, document)
}

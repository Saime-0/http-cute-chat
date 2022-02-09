package healer

import (
	"context"
	"fmt"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"time"
)

type LogRow struct {
	Time  time.Time
	Level string
	Body  interface{}
}

func (h *Healer) Log(document interface{}) {
	if h.Output <= clog.Multiple {
		//b, _ := json.MarshalIndent(document, "", " ")
		//fmt.Println(string(b))
		fmt.Printf("%#v\n", document)
	}
	if h.Output >= clog.Multiple {
		ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
		defer cancel()
		_, err := h.db.Collection("logs").InsertOne(ctx, document)
		h.MonitorLogger(err)
	}
}

func (h *Healer) log(lvl clog.LogLevel, document interface{}) {
	h.Log(&LogRow{
		Time:  time.Now(),
		Level: lvl.String(),
		Body:  document,
	})
}

func (h *Healer) Emergency(document interface{}) {
	h.log(clog.Emergency, document)
}
func (h *Healer) Alert(document interface{}) {
	h.log(clog.Emergency, document)
}
func (h *Healer) Critical(document interface{}) {
	h.log(clog.Critical, document)
}
func (h *Healer) Error(document interface{}) {
	h.log(clog.Error, document)
}
func (h *Healer) Warning(document interface{}) {
	h.log(clog.Warning, document)
}
func (h *Healer) Notice(document interface{}) {
	h.log(clog.Notice, document)
}
func (h *Healer) Info(document interface{}) {
	h.log(clog.Info, document)
}
func (h *Healer) Debug(document interface{}) {
	h.log(clog.Debug, document)
}

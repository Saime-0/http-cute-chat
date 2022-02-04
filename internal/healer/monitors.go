package healer

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/internal/res"
)

func (h Healer) MonitorLogger(err error) (fail bool) {
	fmt.Println("MonitorLogger:", err) // debug
	if err != nil {
		err = h.stateMachine.Indicators[res.IndicatorLogger].SetState(res.FailedDBConnection)
		if err != nil {
			println("MonitorLogger:", err.Error()) // debug
		}
		return true
	}
	return
}

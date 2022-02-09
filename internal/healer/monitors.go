package healer

import (
	"fmt"
	"github.com/saime-0/http-cute-chat/internal/res"
)

func (h Healer) MonitorLogger(err error) {
	if err != nil {
		err = h.stateMachine.Indicators[res.IndicatorLogger].SetState(res.FailedDBConnection)
		if err != nil {
			fmt.Println("MonitorLogger:", err) // debug
		}
	}
}

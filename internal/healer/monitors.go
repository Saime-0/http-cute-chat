package healer

import "github.com/saime-0/http-cute-chat/internal/res"

func (h Healer) MonitorLogger(err error) (fail bool) {
	if err != nil {
		h.stateMachine.Indicators[res.IndicatorLogger].SetState(res.FailedDBConnection)
		return true
	}
	return
}

package healer

import (
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/pkg/fsm"
	"time"
)

func (h *Healer) createLoggerIndicator() (err error) {

	_, err = h.stateMachine.AddIndicator(
		res.IndicatorLogger,
		res.OK,
		fsm.States{
			res.OK:                 h.loggerStateOK(),
			res.FailedDBConnection: h.loggerStateFailedDBConnection(),
			//res.FailedDBConnection:  h.loggerStateFailedDBConnection(),
			//res.RepairingConnection: h.loggerStateRepairingConnection(),
		},
	)

	return err
}

func (h *Healer) loggerStateOK() *fsm.State {
	return fsm.NewState(func(_ *fsm.Indicator) error { // когда возвращается в нормальное состояние
		h.Notice(res.ConnectionToTheLogDBHasBeenSuccessfullyRestored)
		return nil
	})
}

func (h *Healer) loggerStateFailedDBConnection() *fsm.State {
	return fsm.NewState(func(indicator *fsm.Indicator) error { // пропало соединение с бд данных, а возможно и в чем то другом проблема надо проверить

		go func() {
			h.Info(res.StartingLogDBConnectionRecoveryService)

			for i := 0; i < rules.AllowedConnectionShutdownDuration/2; i++ {
				err := h.PingDB()
				if err == nil {
					err = indicator.SetState(res.OK)
					if err != nil {
						panic("not handling")
					}
					return
				}
				time.Sleep(time.Second * 2)
			}

			h.Alert(res.ConnectionToDatabaseCouldNotBeRestored)
		}()
		return nil
	})
}

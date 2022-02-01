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
			res.OK:                  h.loggerStateOK(),
			res.FailedDBConnection:  fsm.NewTransitionalState(),
			res.RepairingConnection: h.loggerStateRepairingConnection(),
		},
	)

	return err
}

func (h *Healer) loggerStateOK() *fsm.State {
	return fsm.NewState(func() error { // когда возвращается в нормальное состояние
		h.services.Logger.Notice(res.ConnectionToTheLogDBHasBeenSuccessfullyRestored)
		return nil
	})
}

func (h *Healer) loggerStateRepairingConnection() *fsm.State {
	return fsm.NewState(func() error { // пропало соединение с бд данных, а возможно и в чем то другом проблема надо проверить
		//_, ok := h.services.Cache.GetState(res.CacheCurrentReconnectionAttemptToLogDB)
		//if ok {
		//	return nil
		//}
		//h.services.Cache.SetState(
		//	res.CacheCurrentReconnectionAttemptToLogDB,
		//	rules.NumberOfAttemptsToConnectToTheLogDBBeforeTheAlert,
		//)
		indicator := h.stateMachine.Indicators[res.IndicatorLogger]
		state := indicator.GetState().(res.LocalKeys)
		if state != res.OK {
			println("state != res.ok; return nil") // debug
			return nil
		}
		//err := indicator.SetState(res.FailedDBConnection)
		//if err != nil {
		//	return errors.Wrap(err, res.FailedToSetTheState)
		//}
		//println("indicator set state") // debug
		h.services.Logger.Warning(res.StartingLogDBConnectionRecoveryService)
		go func() {
			sleep := time.Second * rules.InitialIntervalBetweenAttemptsToReconnectToTheDatabase
			for i := rules.NumberOfAttemptsToConnectToTheLogDBBeforeTheAlert; i >= 0; i -= 1 {
				time.Sleep(sleep)
				println("logger not working, attemptin to repaire - num", i) // debug
				if h.services.Logger.ReconnectToDB() == nil {
					indicator.SetState(res.OK)
					println("logger repaired") // debug
					return
				}
				sleep *= rules.ConnectionIntervalComplexityMultiplier
			}
			h.services.Logger.Alert(res.ConnectionToDatabaseCouldNotBeRestored)
			println("logger repairing failure") // debug
		}()
		return nil
	})
}

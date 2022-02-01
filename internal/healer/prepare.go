package healer

import (
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (h *Healer) prepareHealer() (err error) {

	errHandler := kit.CreateErrorHandler()
	switch {
	case errHandler.Fail(h.createLoggerIndicator()):
		fallthrough

	case false:
		goto handleError
	}

	return nil

handleError:
	return errHandler.Err
}

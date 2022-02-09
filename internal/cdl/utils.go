package cdl

import (
	"github.com/pkg/errors"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (c *parentCategory) getRequest(ptr chanPtr) *baseRequest {
	request, ok := c.Requests[ptr]
	if !ok { // если еще не создавали то надо паниковать
		c.Dataloader.healer.Alert(errors.Wrap(errors.New("c.Requests not exists"), utils.GetCallerPos()))
		panic("c.Requests not exists by" + ptr)
	}
	return request
}

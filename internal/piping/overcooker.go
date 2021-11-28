package piping

import (
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
)

// deprecated
type overcooker struct{}

func (o *overcooker) Chat(before *models.Chat) model.Chat {
	return model.Chat{}
}

package resolver

import (
	"github.com/saime-0/http-cute-chat/internal/cdl"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/service"
	"github.com/saime-0/http-cute-chat/internal/subix"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	Services   *service.Services
	Config     *config.Config2
	Piper      *piper.Pipeline
	Healer     *healer.Healer
	Subix      *subix.Subix
	Dataloader *cdl.Dataloader
}

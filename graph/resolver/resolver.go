package resolver

import (
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	Services *service.Services
	Config   *config.Config
	Piper    *piper.Pipeline
}

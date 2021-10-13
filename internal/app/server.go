package app

import (
	"context"
	"net/http"
)

func NewHttpServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}

}

func (a *ApiServer) Run() error {
	return a.httpServer.ListenAndServe()
}

func (a *ApiServer) Stop(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}

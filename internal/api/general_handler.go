package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	v1 "github.com/saime-0/http-cute-chat/internal/api/v1"
	"github.com/saime-0/http-cute-chat/internal/service"
)

type GeneralHandler struct {
	Services *service.Services
}

func NewGeneralHandler(services *service.Services) *GeneralHandler {
	return &GeneralHandler{
		Services: services,
	}
}

func (h *GeneralHandler) Init() *mux.Router {
	r := mux.NewRouter()
	r.Use(logRequest)
	h.initAPI(r)
	return r
}

func (h *GeneralHandler) initAPI(router *mux.Router) {
	handlerV1 := v1.NewHandler(h.Services)
	api := router.PathPrefix("/api/").Subrouter()
	handlerV1.Init(api)

}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf(
			"completed with %d %s in %v\n",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

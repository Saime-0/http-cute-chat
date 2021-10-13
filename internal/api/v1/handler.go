package v1

import (
	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/service"
)

type Handler struct {
	Services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) Init(r *mux.Router) {
	v1 := r.PathPrefix("/v1/").Subrouter()
	{
		h.initUsersRoutes(v1)
		// h.initUsersRoutes(v1)
		// h.initCoursesRoutes(v1)
		// h.initStudentsRoutes(v1)
		// h.initCallbackRoutes(v1)
		// h.initAdminRoutes(v1)

		// v1.GET("/settings", h.setSchoolFromRequest, h.getSchoolSettings)
		// v1.GET("/promocodes/:code", h.setSchoolFromRequest, h.getPromo)
		// v1.GET("/offers/:id", h.setSchoolFromRequest, h.getOffer)
	}
}

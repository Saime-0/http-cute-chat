package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) initDialogsRoutes(r *mux.Router) {
	dialogs := r.PathPrefix("/dialogs/").Subrouter()
	{
		//

		authenticated := dialogs.PathPrefix("/").Subrouter()
		authenticated.Use(h.checkAuth)
		{
			// POST
			authenticated.HandleFunc("/{user-id}/messages/", h.SendMessageToUser).Methods(http.MethodPost)
			// GET
			authenticated.HandleFunc("/", h.GetListCompanions).Methods(http.MethodGet)
			authenticated.HandleFunc("/{user-id}/messages/", h.GetListDialogMessages).Methods(http.MethodGet)
			authenticated.HandleFunc("/dialogs/{user-id}/messages/{message-id}/", h.GetDialogMessage).Methods(http.MethodGet)
			// PUT
		}
	}
}

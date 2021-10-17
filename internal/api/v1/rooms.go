package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) initRoomsRoutes(r *mux.Router) {
	rooms := r.PathPrefix("/chats/").Subrouter()
	{
		//

		authenticated := rooms.PathPrefix("/").Subrouter()
		authenticated.Use(h.checkAuth)
		{
			// POST
			authenticated.HandleFunc("/{chat-id}/rooms/{room-id}/messages/", h.SendMessageToRoom).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/rooms/", h.CreateRoom).Methods(http.MethodPost)
			// GET
			authenticated.HandleFunc("/{chat-id}/rooms/{room-id}/messages/", h.GetListRoomMessages).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/rooms/{room-id}/messages/{message-id}/", h.GetRoomMessage).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{chat-id}/rooms/{room-id}/data/", h.SetChatData).Methods(http.MethodPut)
		}
	}
}

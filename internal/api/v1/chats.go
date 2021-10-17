package v1

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) initChatsRoutes(r *mux.Router) {
	chats := r.PathPrefix("/chats/").Subrouter()
	{
		// GET
		chats.HandleFunc("/d/{chat-id}/", h.GetChatByDomain).Methods(http.MethodGet)
		chats.HandleFunc("/{chat-id}/", h.GetChatByID).Methods(http.MethodGet)
		chats.HandleFunc("/", h.GetChatsByName).Methods(http.MethodGet)

		authenticated := chats.PathPrefix("/").Subrouter()
		authenticated.Use(h.checkAuth)
		{
			// POST
			authenticated.HandleFunc("/", h.CreateChat).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/join/", h.AddUserToChat).Methods(http.MethodPost)
			// GET
			authenticated.HandleFunc("/{chat-id}/data/", h.GetChatData).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/members/", h.GetChatMembers).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/rooms/", h.GetChatRooms).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{chat-id}/data/", h.SetChatData).Methods(http.MethodPut)
		}
	}
}

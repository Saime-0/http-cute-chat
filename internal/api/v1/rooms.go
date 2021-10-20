package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (h *Handler) initRoomsRoutes(r *mux.Router) {
	rooms := r.PathPrefix("/rooms/").Subrouter()
	{
		//

		authenticated := rooms.PathPrefix("/").Subrouter()
		authenticated.Use(h.checkAuth)
		{
			// POST
			authenticated.HandleFunc("/{room-id}/messages/", h.SendMessageToRoom).Methods(http.MethodPost)
			authenticated.HandleFunc("/", h.CreateRoom).Methods(http.MethodPost)
			// GET
			authenticated.HandleFunc("/{room-id}/messages/", h.GetListRoomMessages).Methods(http.MethodGet)
			authenticated.HandleFunc("/{room-id}/messages/{message-id}/", h.GetRoomMessage).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{room-id}/data/", h.SetChatData).Methods(http.MethodPut)
		}
	}
}

func (h *Handler) SendMessageToRoom(w http.ResponseWriter, r *http.Request) {
	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		panic(err)
	}
	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		panic(err)
	}
	message_id, err := h.Services.Repos.Rooms.CreateMessage(room_id, message)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, &models.MessageID{message_id})
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		panic(err)
	}
	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		panic(err)
	}
	message_id, err := h.Services.Repos.Rooms.CreateMessage(room_id, message)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, &models.MessageID{message_id})
}

package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
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
			// GET
			authenticated.HandleFunc("/{room-id}/messages/", h.GetRoomMessages).Methods(http.MethodGet)
			authenticated.HandleFunc("/{room-id}/messages/{message-id}/", h.GetRoomMessage).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{room-id}/data/", h.UpdateRoomData).Methods(http.MethodPut)
		}
	}
}

func (h *Handler) SendMessageToRoom(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
	if err != nil {
		panic(err)
	}
	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		panic(err)
	}
	chat_id, err := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if err != nil {
		panic(err)
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		panic(err)
	}
	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		panic(err)
	}
	message.Author = user_id
	message_id, err := h.Services.Repos.Rooms.CreateMessage(room_id, message)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, &models.MessageID{ID: message_id})
}

func (h *Handler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
	if err != nil {
		panic(err)
	}
	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		panic(err)
	}
	chat_id, err := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if err != nil {
		panic(err)
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		panic(err)
	}
	message_list, err := h.Services.Repos.Rooms.GetMessages(room_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, message_list)
}

func (h *Handler) GetRoomMessage(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
	if err != nil {
		panic(err)
	}
	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		panic(err)
	}
	message_id, err := strconv.Atoi(mux.Vars(r)["message-id"])
	if err != nil {
		panic(err)
	}
	chat_id, err := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if err != nil {
		panic(err)
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		panic(err)
	}
	message, err := h.Services.Repos.Rooms.GetMessageInfo(message_id, room_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, message)
}

func (h *Handler) UpdateRoomData(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
	if err != nil {
		panic(err)
	}
	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		panic(err)
	}
	chat_id, err := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if err != nil {
		panic(err)
	}
	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		panic(err)
	}
	room_data := &models.UpdateRoomData{}
	err = json.NewDecoder(r.Body).Decode(&room_data)
	if err != nil {
		panic(err)
	}
	err = h.Services.Repos.Rooms.UpdateRoomData(room_id, room_data)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, "")
}

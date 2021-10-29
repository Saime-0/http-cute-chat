package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saime-0/http-cute-chat/internal/api/rules"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (h *Handler) initRoomsRoutes(r *mux.Router) {
	authenticated := r.PathPrefix("/rooms").Subrouter()
	authenticated.Use(checkAuth)
	{
		// POST
		authenticated.HandleFunc("/{room-id}/messages", h.SendMessageToRoom).Methods(http.MethodPost)
		// GET
		authenticated.HandleFunc("/{room-id}/messages", h.GetRoomMessages).Methods(http.MethodGet)
		authenticated.HandleFunc("/{room-id}/messages/{message-id}", h.GetRoomMessage).Methods(http.MethodGet)
		// PUT
		authenticated.HandleFunc("/{room-id}/data", h.UpdateRoomData).Methods(http.MethodPut)
	}
}

func (h *Handler) SendMessageToRoom(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(room_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chat_id, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	message.Author = user_id
	message_id, err := h.Services.Repos.Messages.CreateMessageInRoom(room_id, message)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.MessageID{ID: message_id})
}

func (h *Handler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	offset, err := parseOffsetFromQuery(w, r)
	if err != nil {

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(room_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chat_id, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) &&
		!h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	message_list, err := h.Services.Repos.Messages.GetMessagesFromRoom(room_id, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, message_list)
}

func (h *Handler) GetRoomMessage(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(room_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chat_id, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) &&
		!h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	message_id, err := strconv.Atoi(mux.Vars(r)["message-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	message, err := h.Services.Repos.Messages.GetMessageFromRoom(message_id, room_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, message)
}

func (h *Handler) UpdateRoomData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	room_id, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	chat_id, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room_id)
	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) &&
		!h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room_data := &models.UpdateRoomData{}
	err = json.NewDecoder(r.Body).Decode(&room_data)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	err = h.Services.Repos.Rooms.UpdateRoomData(room_id, room_data)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saime-0/http-cute-chat/internal/api/rules"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (h *Handler) initChatsRoutes(r *mux.Router) {
	chats := r.PathPrefix("/chats").Subrouter()
	{
		// ! /chats/data/ equal /chats/{chat-id}/
		authenticated := chats.PathPrefix("").Subrouter()
		authenticated.Use(checkAuth)
		{
			// POST
			authenticated.HandleFunc("", h.CreateChat).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/join", h.AddUserToChat).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/rooms", h.CreateRoom).Methods(http.MethodPost)
			// GET
			authenticated.HandleFunc("/{chat-id}/data", h.GetChatData).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/members", h.GetChatMembers).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/rooms", h.GetChatRooms).Methods(http.MethodGet)
			authenticated.HandleFunc("/owned", h.GetUserOwnedChats).Methods(http.MethodGet)
			authenticated.HandleFunc("/involved", h.GetUserChats).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{chat-id}/data", h.UpdateChatData).Methods(http.MethodPut)
		}
		// GET
		chats.HandleFunc("/d/{chat-domain}", h.GetChatByDomain).Methods(http.MethodGet)
		chats.HandleFunc("/{chat-id}", h.GetChatByID).Methods(http.MethodGet)
		chats.HandleFunc("", h.GetChatsByName).Methods(http.MethodGet)
	}
}

func (h *Handler) GetChatByDomain(w http.ResponseWriter, r *http.Request) {
	chat_domain := mux.Vars(r)["chat-domain"]
	if !validateDomain(chat_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByDomain(chat_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatInfoByDomain(chat_domain)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatInfoByID(chat_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatsByName(w http.ResponseWriter, r *http.Request) {
	name_fragment := r.URL.Query().Get("name")
	if len(name_fragment) > rules.NameMaxLength || len(name_fragment) == 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil && r.URL.Query().Get("offset") != "" {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if offset < 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}

	chat_list, err := h.Services.Repos.Chats.GetChatsByNameFragment(name_fragment, offset)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat := &models.CreateChat{}
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	count_chats, err := h.Services.Repos.Users.GetCountUserOwnedChats(user_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if count_chats >= rules.MaxCountOwnedChats {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	if h.Services.Repos.Units.IsUnitExistsByDomain(chat.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	if !validateDomain(chat.Domain) || !validateName(chat.Name) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	chat_id, err := h.Services.Repos.Chats.CreateChat(user_id, chat)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, &models.ChatID{ID: chat_id})

}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room := &models.CreateRoom{}
	err = json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	count_rooms, err := h.Services.Repos.Chats.GetCountRooms(chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if count_rooms >= rules.MaxCountRooms {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	parent_chat, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room.ParentID)
	if parent_chat != chat_id {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	room_id, err := h.Services.Repos.Rooms.CreateRoom(chat_id, room)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, &models.RoomID{ID: room_id})
}

func (h *Handler) AddUserToChat(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	count_chats, err := h.Services.Repos.Chats.GetCountUserChats(user_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if count_chats >= rules.MaxUserChats {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	err = h.Services.Repos.Chats.AddUserToChat(user_id, chat_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	chat_data, err := h.Services.Repos.Chats.GetChatDataByID(chat_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, chat_data)
}

func (h *Handler) GetChatMembers(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	user_list, err := h.Services.Repos.Chats.GetChatMembers(chat_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) GetChatRooms(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room_list, err := h.Services.Repos.Rooms.GetChatRooms(chat_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, room_list)
}

func (h *Handler) GetUserOwnedChats(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_list, err := h.Services.Repos.Chats.GetChatsOwnedUser(user_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_list, err := h.Services.Repos.Chats.GetChatsInvolvedUser(user_id)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) UpdateChatData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)
		return
	}

	if !h.Services.Repos.Units.IsUnitExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	chat_data := &models.UpdateChatData{}
	err = json.NewDecoder(r.Body).Decode(&chat_data)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	err = h.Services.Repos.Chats.UpdateChatData(chat_id, chat_data)
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		return
	}

	responder.Respond(w, http.StatusOK, nil)
}

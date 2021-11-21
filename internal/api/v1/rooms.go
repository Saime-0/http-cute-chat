package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/saime-0/http-cute-chat/internal/api/validator"

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
		authenticated.HandleFunc("/{room-id}/form", h.GetRoomForm).Methods(http.MethodGet)
		// PUT
		authenticated.HandleFunc("/{room-id}/data", h.UpdateRoomData).Methods(http.MethodPut)
		authenticated.HandleFunc("/{room-id}/form", h.SetRoomForm).Methods(http.MethodPut)
		// DELETE
		authenticated.HandleFunc("/{room-id}/form", h.ClearRoomForm).Methods(http.MethodDelete)
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

	room, err := h.Services.Repos.Rooms.GetRoom(room_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	//- private and !room_id and (!manage or !nil)
	//- !(room.Private && (role.ManageRooms && role.RoomID == 0 || role.RoomID == room_id) || !room.Private || h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) && )
	//if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) && room.Private && role.RoomID != room_id && (!role.ManageRooms || role.RoomID != 0) {
	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) || role.RoomID == room_id || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	message.Author = user_id

	msg_type := rules.UserMsg
	if h.Services.Repos.Rooms.RoomFormIsSet(room_id) {
		msg_type = rules.FormattedMsg
		var input_choice models.FormCompleted
		err := json.Unmarshal([]byte(message.Body), &input_choice)
		if err != nil {
			responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

			return
		}
		room_form, err := h.Services.Repos.Rooms.GetRoomForm(room_id)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

			panic(err)
		}
		choice, aerr := MatchMessageType(&input_choice, &room_form)
		if aerr != nil {
			responder.Error(w, http.StatusBadRequest, aerr)

			return
		}
		msg_body, err := json.Marshal(choice)
		if err != nil {
			responder.Error(w, http.StatusBadRequest, rules.ErrDataRetrieved)

		}
		message.Body = string(msg_body)
	}

	message_id, err := h.Services.Repos.Messages.CreateMessageInRoom(room_id, msg_type, message)
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

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

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

	room, err := h.Services.Repos.Rooms.GetRoom(room_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) || role.RoomID == room_id || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	// todo: member have permissions
	member, err := h.Services.Repos.Chats.GetMemberInfo(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrDataRetrieved)

		panic(err)
	}
	message_list, err := h.Services.Repos.Messages.GetMessagesFromRoom(room_id, member.JoinedAt, offset)
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

	room, err := h.Services.Repos.Rooms.GetRoom(room_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	role, err := h.Services.Repos.Chats.GetUserRoleData(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) || role.RoomID == room_id || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

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

func (h *Handler) GetRoomForm(w http.ResponseWriter, r *http.Request) {
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

	room, err := h.Services.Repos.Rooms.GetRoom(room_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	role, err := h.Services.Repos.Chats.GetUserRoleData(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) || role.RoomID == room_id || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	if !h.Services.Repos.Rooms.RoomFormIsSet(room_id) {
		responder.Respond(w, http.StatusOK, nil)

		return
	}

	form, err := h.Services.Repos.Rooms.GetRoomForm(room_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, form)
}

func (h *Handler) SetRoomForm(w http.ResponseWriter, r *http.Request) {
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

	role, err := h.Services.Repos.Chats.GetUserRoleData(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) || role.RoomID == room_id || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	form_pattern := &models.FormPattern{}
	err = json.NewDecoder(r.Body).Decode(&form_pattern)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}
	if !validator.ValidateRoomForm(form_pattern) {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	format, err := io.ReadAll(r.Body)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)

		panic(err)
	}
	err = h.Services.Repos.Rooms.UpdateRoomForm(room_id, string(format))
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) ClearRoomForm(w http.ResponseWriter, r *http.Request) {
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

	role, err := h.Services.Repos.Chats.GetUserRoleData(user_id, chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) || role.RoomID == room_id || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	err = h.Services.Repos.Rooms.UpdateRoomForm(room_id, "")
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

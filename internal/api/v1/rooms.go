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
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room, err := h.Services.Repos.Rooms.Room(roomId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	//- private and !room_id and (!manage or !nil)
	//- !(room.Private && (role.ManageRooms && role.RoomID == 0 || role.RoomID == room_id) || !room.Private || h.Services.repos.Chats.UserIsChatOwner(user_id, chat_id) && )
	//if !h.Services.repos.Chats.UserIsChatOwner(user_id, chat_id) && room.Private && role.RoomID != room_id && (!role.ManageRooms || role.RoomID != 0) {
	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.RoomID == roomId || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	message.Author = userId

	msgType := rules.UserMsg
	if h.Services.Repos.Rooms.FormIsSet(roomId) {
		msgType = rules.FormattedMsg
		var inputChoice models.FormCompleted
		err := json.Unmarshal([]byte(message.Body), &inputChoice)
		if err != nil {
			responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

			return
		}
		roomForm, err := h.Services.Repos.Rooms.RoomForm(roomId)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

			panic(err)
		}
		choice, aerr := MatchMessageType(&inputChoice, &roomForm)
		if aerr != nil {
			responder.Error(w, http.StatusBadRequest, aerr)

			return
		}
		msgBody, err := json.Marshal(choice)
		if err != nil {
			responder.Error(w, http.StatusBadRequest, rules.ErrDataRetrieved)

		}
		message.Body = string(msgBody)
	}

	messageId, err := h.Services.Repos.Messages.CreateMessageInRoom(roomId, msgType, message)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.MessageID{ID: messageId})
}

func (h *Handler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) &&
		!h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room, err := h.Services.Repos.Rooms.Room(roomId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.RoomID == roomId || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	// todo: member have permissions
	member, err := h.Services.Repos.Chats.GetMemberInfo(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrDataRetrieved)

		panic(err)
	}
	messageList, err := h.Services.Repos.Messages.GetMessagesFromRoom(roomId, member.JoinedAt, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, messageList)
}

func (h *Handler) GetRoomMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) &&
		!h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room, err := h.Services.Repos.Rooms.Room(roomId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.RoomID == roomId || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	messageId, err := strconv.Atoi(mux.Vars(r)["message-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	message, err := h.Services.Repos.Messages.Message(messageId, roomId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, message)
}

func (h *Handler) UpdateRoomData(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) &&
		!h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	roomData := &models.UpdateRoomData{}
	err = json.NewDecoder(r.Body).Decode(&roomData)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	err = h.Services.Repos.Rooms.UpdateRoomData(roomId, roomData)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetRoomForm(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room, err := h.Services.Repos.Rooms.Room(roomId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(!room.Private || h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.RoomID == roomId || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	if !h.Services.Repos.Rooms.FormIsSet(roomId) {
		responder.Respond(w, http.StatusOK, nil)

		return
	}

	form, err := h.Services.Repos.Rooms.RoomForm(roomId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, form)
}

func (h *Handler) SetRoomForm(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.RoomID == roomId || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	formPattern := &models.FormPattern{}
	err = json.NewDecoder(r.Body).Decode(&formPattern)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}
	if !validator.ValidateRoomForm(formPattern) {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}
	format, err := io.ReadAll(r.Body)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)

		panic(err)
	}
	err = h.Services.Repos.Rooms.UpdateRoomForm(roomId, string(format))
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) ClearRoomForm(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	roomId, err := strconv.Atoi(mux.Vars(r)["room-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Rooms.RoomExistsByID(roomId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrRoomNotFound)

		return
	}

	chatId, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(roomId)
	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.RoomID == roomId || role.ManageRooms && role.RoomID == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrPrivateRoom)

		return
	}

	err = h.Services.Repos.Rooms.UpdateRoomForm(roomId, "")
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

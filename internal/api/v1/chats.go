package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/api/validator"

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
			authenticated.HandleFunc("/{chat-id}/links", h.CreateInviteLink).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/banlist/{user-id}", h.BanUserInChat).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/roles", h.CreateRole).Methods(http.MethodPost)
			authenticated.HandleFunc("/{chat-id}/members/{user-id}/role", h.AddRoleToUser).Methods(http.MethodPost) // ! переработать
			// GET
			authenticated.HandleFunc("/{chat-id}/data", h.GetChatData).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/members", h.GetChatMembers).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/rooms", h.GetChatRooms).Methods(http.MethodGet)
			authenticated.HandleFunc("/owned", h.GetUserOwnedChats).Methods(http.MethodGet)
			authenticated.HandleFunc("/involved", h.GetUserChats).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/links", h.GetInviteLinks).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/banlist", h.GetChatBanlist).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/roles", h.GetChatRoles).Methods(http.MethodGet)
			// authenticated.HandleFunc("/{chat-id}/roles/data", h.GetChatRolesData).Methods(http.MethodGet)
			authenticated.HandleFunc("/{chat-id}/members/{user-id}/role", h.GetUserRole).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{chat-id}/data", h.UpdateChatData).Methods(http.MethodPut)
			authenticated.HandleFunc("/{chat-id}/roles/{role-id}", h.UpdateRoleData).Methods(http.MethodPut)
			// DELETE
			authenticated.HandleFunc("/{chat-id}/leave", h.RemoveUserFromChat).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/links/{invite-code}", h.DeleteInviteLink).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/banlist/{user-id}", h.UnbanUserInChat).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/roles/{role-id}", h.RemoveChatRole).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/members/{user-id}/role", h.TakeUserRole).Methods(http.MethodDelete)

		}
		// GET
		chats.HandleFunc("/d/{chat-domain}", h.GetChatByDomain).Methods(http.MethodGet)
		chats.HandleFunc("/{chat-id}", h.GetChatByID).Methods(http.MethodGet)
		chats.HandleFunc("", h.GetChatsByName).Methods(http.MethodGet)
	}
}

func (h *Handler) GetChatByDomain(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatDomain := mux.Vars(r)["chat-domain"]
	if !validator.ValidateDomain(chatDomain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByDomain(chatDomain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatByDomain(chatDomain)
	finalInspectionDatabase(w, err)
	if h.Services.Repos.Chats.ChatIsPrivate(chat.ID) && !h.Services.Repos.Chats.UserIsChatMember(userId, chat.ID) {
		responder.Error(w, http.StatusForbidden, rules.ErrNoAccess)

		return
	}

	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	if h.Services.Repos.Chats.ChatIsPrivate(chatId) && !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusForbidden, rules.ErrNoAccess)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatByID(chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatsByName(w http.ResponseWriter, r *http.Request) {
	nameFragment := r.URL.Query().Get("name")
	if len(nameFragment) > rules.NameMaxLength || len(nameFragment) == 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	chatList, err := h.Services.Repos.Chats.GetChatsByNameFragment(nameFragment, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chatList)
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chat := &models.CreateChat{}
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	countChats, err := h.Services.Repos.Users.GetCountUserOwnedChats(userId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if countChats >= rules.MaxCountOwnedChats {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	if h.Services.Repos.Chats.ChatExistsByDomain(chat.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	if !validator.ValidateDomain(chat.Domain) || !validator.ValidateName(chat.Name) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	chatId, err := h.Services.Repos.Chats.CreateChat(userId, chat)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.ChatID{ID: chatId})

}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}
	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageRooms) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room := &models.CreateRoom{}
	err = json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	countRooms, err := h.Services.Repos.Chats.GetCountRooms(chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if countRooms >= rules.MaxCountRooms {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	if h.Services.Repos.Rooms.RoomExistsByID(room.ParentID) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	parentChat, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room.ParentID)
	if parentChat != chatId {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	roomId, err := h.Services.Repos.Rooms.CreateRoom(chatId, room)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.RoomID{ID: roomId})
}

func (h *Handler) AddUserToChat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if h.Services.Repos.Chats.ChatIsPrivate(chatId) {
		responder.Error(w, http.StatusForbidden, rules.ErrNoAccess)

		return
	}

	countChats, err := h.Services.Repos.Chats.GetCountUserChats(userId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if countChats >= rules.MaxUserChats {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	countMembers, err := h.Services.Repos.Chats.CountMembers(chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if countMembers >= rules.MaxMembersOnChat {
		responder.Error(w, http.StatusBadRequest, rules.ErrMembersLimitHasBeenReached)

		return
	}

	err = h.Services.Repos.Chats.AddUserToChat(userId, chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) RemoveUserFromChat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	err = h.Services.Repos.Chats.RemoveUserFromChat(userId, chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatData(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageChat) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	chatData, err := h.Services.Repos.Chats.GetChatDataByID(chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chatData)
}

func (h *Handler) GetChatMembers(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) &&
		!h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	userList, err := h.Services.Repos.Chats.Members(chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, userList)
}

func (h *Handler) GetChatRooms(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) &&
		!h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	roomList, err := h.Services.Repos.Rooms.GetChatRooms(chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, roomList)
}

func (h *Handler) GetUserOwnedChats(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	chatList, err := h.Services.Repos.Chats.GetChatsOwnedUser(userId, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chatList)
}

func (h *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	chatList, err := h.Services.Repos.Chats.GetChatsInvolvedUser(userId, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chatList)
}

func (h *Handler) UpdateChatData(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)
		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageChat) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	chatData := &models.UpdateChatData{}
	err = json.NewDecoder(r.Body).Decode(&chatData)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if h.Services.Repos.Chats.ChatExistsByDomain(chatData.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	err = h.Services.Repos.Chats.UpdateChatData(chatId, chatData)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) CreateInviteLink(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	// мб эту часть кода свернуть в функцию, то есть то надо взять ид и знать что существует такой юнит или комната
	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}
	// до сюда

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageChat) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	inputLink := &models.InviteInput{}
	err = json.NewDecoder(r.Body).Decode(&inputLink)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if inputLink.LifeTime != 0 && !validator.ValidateLifetime(inputLink.LifeTime) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}
	if inputLink.Aliens != 0 && !validator.ValidateAliens(inputLink.Aliens) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}

	count, err := h.Services.Repos.Chats.GetCountLinks(chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if count > rules.MaxInviteLinks {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	link, err := h.Services.Repos.Chats.CreateInviteLink(
		&models.CreateInvite{
			ChatID: chatId,
			Aliens: inputLink.Aliens,
			Exp:    inputLink.LifeTime + time.Now().UTC().Unix(),
		},
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, link)
}

func (h *Handler) GetInviteLinks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageChat) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	links, err := h.Services.Repos.Chats.Invites(chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, links)
}

func (h *Handler) DeleteInviteLink(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	linkCode := mux.Vars(r)["invite-code"]
	if !h.Services.Repos.Chats.InviteIsRelevant(linkCode) {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	err = h.Services.Repos.Chats.DeleteInviteLinkByCode(linkCode)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) BanUserInChat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)
	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(targetId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageMembers) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	role, err = h.Services.Repos.Chats.GetUserRoleData(targetId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if h.Services.Repos.Chats.UserIsChatOwner(targetId, chatId) || role.ManageRooms {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Chats.AddToBanlist(targetId, chatId)
	finalInspectionDatabase(w, err)
	err = h.Services.Repos.Chats.RemoveUserFromChat(targetId, chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) UnbanUserInChat(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageMembers) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsBannedInChat(targetId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Chats.RemoveFromBanlist(targetId, chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatBanlist(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageChat) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	users, err := h.Services.Repos.Chats.Banlist(chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, users)
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageRooms) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	countRoles, err := h.Services.Repos.Chats.GetCountChatRoles(chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if countRoles > rules.MaxRolesInChat {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	roleModel := &models.CreateRole{}
	err = json.NewDecoder(r.Body).Decode(&roleModel)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	roleId, err := h.Services.Repos.Chats.CreateRoleInChat(chatId, roleModel)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, models.RoleID{ID: roleId})
}

func (h *Handler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageRooms) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(targetId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

	err = h.Services.Repos.Chats.GiveRole(targetId, chatId) // ! fix role_id?
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatRoles(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	switch {
	case h.Services.Repos.Chats.UserIsChatOwner(userId, chatId):
		roles, err := h.Services.Repos.Chats.GetChatRolesData(chatId)
		finalInspectionDatabase(w, err)

		responder.Respond(w, http.StatusOK, roles)

	case h.Services.Repos.Chats.UserIsChatMember(userId, chatId):
		roles, err := h.Services.Repos.Chats.Roles(chatId)
		finalInspectionDatabase(w, err)

		responder.Respond(w, http.StatusOK, roles)
	default:
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}
}

func (h *Handler) GetUserRole(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(userId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}
	if !h.Services.Repos.Chats.UserIsChatMember(targetId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(targetId, chatId)
	finalInspectionDatabase(w, err)

	switch {
	case h.Services.Repos.Chats.UserIsChatOwner(userId, chatId):
		responder.Respond(w, http.StatusOK, role)
	default:
		if !role.Visible {
			responder.Error(w, http.StatusBadRequest, rules.ErrRoleHidden)

			return
		}
		responder.Respond(w, http.StatusOK, models.RoleInfo{
			ID:       role.ID,
			RoleName: role.RoleName,
			Color:    role.Color,
		})
	}

}

func (h *Handler) UpdateRoleData(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageRoles) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	roleId, err := strconv.Atoi(mux.Vars(r)["role-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.RoleExistsByID(roleId, chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	updateModel := &models.UpdateRole{}
	err = json.NewDecoder(r.Body).Decode(&updateModel)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	err = h.Services.Repos.Chats.UpdateRoleData(roleId, updateModel)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) RemoveChatRole(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageChat) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	roleId, err := strconv.Atoi(mux.Vars(r)["role-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.RoleExistsByID(roleId, chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	err = h.Services.Repos.Chats.DeleteRole(roleId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) TakeUserRole(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	chatId, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chatId) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	role, err := h.Services.Repos.Chats.GetUserRoleData(userId, chatId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if !(h.Services.Repos.Chats.UserIsChatOwner(userId, chatId) || role.ManageRoles) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}
	if !h.Services.Repos.Chats.UserIsChatMember(targetId, chatId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

	err = h.Services.Repos.Chats.TakeRole(targetId, chatId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
			authenticated.HandleFunc("/{chat-id}/members/{user-id}/role", h.GetUserRole).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{chat-id}/data", h.UpdateChatData).Methods(http.MethodPut)
			authenticated.HandleFunc("/{chat-id}/roles/{role-id}", h.UpdateRoleData).Methods(http.MethodPut)
			// DELETE
			authenticated.HandleFunc("/{chat-id}/leave", h.RemoveUserFromChat).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/links/{invite-code}", h.DeleteInviteLink).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/banlist/{user-id}", h.UnbanUserInChat).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/roles/{role-id}", h.RemoveChatRole).Methods(http.MethodDelete)
			authenticated.HandleFunc("/{chat-id}/roles/{role-id}", h.TakeUserRole).Methods(http.MethodDelete)

		}
		// GET
		chats.HandleFunc("/d/{chat-domain}", h.GetChatByDomain).Methods(http.MethodGet)
		chats.HandleFunc("/{chat-id}", h.GetChatByID).Methods(http.MethodGet)
		chats.HandleFunc("", h.GetChatsByName).Methods(http.MethodGet)
	}
}

func (h *Handler) GetChatByDomain(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_domain := mux.Vars(r)["chat-domain"]
	if !validateDomain(chat_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByDomain(chat_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatByDomain(chat_domain)
	finalInspectionDatabase(w, err)
	if h.Services.Repos.Chats.ChatIsPrivate(chat.ID) && !h.Services.Repos.Chats.UserIsChatMember(user_id, chat.ID) {
		responder.Error(w, http.StatusForbidden, rules.ErrNoAccess)

		return
	}

	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrChatNotFound)

		return
	}

	if h.Services.Repos.Chats.ChatIsPrivate(chat_id) && !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		responder.Error(w, http.StatusForbidden, rules.ErrNoAccess)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatByID(chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatsByName(w http.ResponseWriter, r *http.Request) {
	name_fragment := r.URL.Query().Get("name")
	if len(name_fragment) > rules.NameMaxLength || len(name_fragment) == 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	chat_list, err := h.Services.Repos.Chats.GetChatsByNameFragment(name_fragment, offset)
	finalInspectionDatabase(w, err)

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

	if h.Services.Repos.Chats.ChatExistsByDomain(chat.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	if !validateDomain(chat.Domain) || !validateName(chat.Name) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	chat_id, err := h.Services.Repos.Chats.CreateChat(user_id, chat)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.ChatID{ID: chat_id})

}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
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

	if h.Services.Repos.Rooms.RoomExistsByID(room.ParentID) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	parent_chat, _ := h.Services.Repos.Rooms.GetChatIDByRoomID(room.ParentID)
	if parent_chat != chat_id {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	room_id, err := h.Services.Repos.Rooms.CreateRoom(chat_id, room)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.RoomID{ID: room_id})
}

func (h *Handler) AddUserToChat(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if h.Services.Repos.Chats.ChatIsPrivate(chat_id) {
		responder.Error(w, http.StatusForbidden, rules.ErrNoAccess)

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

	count_members, err := h.Services.Repos.Chats.GetCountChatMembers(chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if count_members >= rules.MaxMembersOnChat {
		responder.Error(w, http.StatusBadRequest, rules.ErrMembersLimitHasBeenReached)

		return
	}

	err = h.Services.Repos.Chats.AddUserToChat(user_id, chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) RemoveUserFromChat(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	err = h.Services.Repos.Chats.RemoveUserFromChat(user_id, chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	chat_data, err := h.Services.Repos.Chats.GetChatDataByID(chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat_data)
}

func (h *Handler) GetChatMembers(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) &&
		!h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	user_list, err := h.Services.Repos.Chats.GetChatMembers(chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) GetChatRooms(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) &&
		!h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	room_list, err := h.Services.Repos.Rooms.GetChatRooms(chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, room_list)
}

func (h *Handler) GetUserOwnedChats(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	chat_list, err := h.Services.Repos.Chats.GetChatsOwnedUser(user_id, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	chat_list, err := h.Services.Repos.Chats.GetChatsInvolvedUser(user_id, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) UpdateChatData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)
		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
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

	if h.Services.Repos.Chats.ChatExistsByDomain(chat_data.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	err = h.Services.Repos.Chats.UpdateChatData(chat_id, chat_data)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) CreateInviteLink(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	// мб эту часть кода свернуть в функцию, то есть то надо взять ид и знать что существует такой юнит или комната
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}
	// до сюда

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	input_link := &models.InviteLinkInput{}
	err = json.NewDecoder(r.Body).Decode(&input_link)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if input_link.LifeTime != 0 && !validateLifetime(input_link.LifeTime) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}
	if input_link.Aliens != 0 && !validateAliens(input_link.Aliens) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}

	count, err := h.Services.Repos.Chats.GetCountLinks(chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if count > rules.MaxInviteLinks {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	link, err := h.Services.Repos.Chats.CreateInviteLink(
		&models.CreateInviteLink{
			ChatID: chat_id,
			Aliens: input_link.Aliens,
			Exp:    input_link.LifeTime + time.Now().UTC().Unix(),
		},
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, link)
}

func (h *Handler) GetInviteLinks(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	links, err := h.Services.Repos.Chats.GetChatLinks(chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, links)
}

func (h *Handler) DeleteInviteLink(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	link_code := mux.Vars(r)["invite-code"]
	if !h.Services.Repos.Chats.InviteLinkIsRelevant(link_code) {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	err = h.Services.Repos.Chats.DeleteInviteLinkByCode(link_code)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) BanUserInChat(w http.ResponseWriter, r *http.Request) {
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	user_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(user_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	if h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

	err = h.Services.Repos.Chats.BanUserInChat(user_id, chat_id)
	finalInspectionDatabase(w, err)
	err = h.Services.Repos.Chats.RemoveUserFromChat(user_id, chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) UnbanUserInChat(w http.ResponseWriter, r *http.Request) {
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	user_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(user_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsBannedInChat(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Chats.UnbanUserInChat(user_id, chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatBanlist(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	users, err := h.Services.Repos.Chats.GetChatBanlist(chat_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, users)
}

func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	count_roles, err := h.Services.Repos.Chats.GetCountChatRoles(chat_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if count_roles > rules.MaxRolesInChat {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	role_model := &models.CreateRole{}
	err = json.NewDecoder(r.Body).Decode(&role_model)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	role_id, err := h.Services.Repos.Chats.CreateRoleInChat(chat_id, role_model)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, models.RoleID{ID: role_id})
}

func (h *Handler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Chats.ChatExistsByID(chat_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrChatNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	member_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(member_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	if !h.Services.Repos.Chats.UserIsChatMember(member_id, chat_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserIsNotChatMember)

		return
	}

}

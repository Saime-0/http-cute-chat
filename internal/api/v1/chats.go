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
			authenticated.HandleFunc("/owned/", h.GetUserOwnedChats).Methods(http.MethodGet)
			authenticated.HandleFunc("/involved/", h.GetUserChats).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/{chat-id}/data/", h.UpdateChatData).Methods(http.MethodPut)
		}
	}
}

func (h *Handler) GetChatByDomain(w http.ResponseWriter, r *http.Request) {
	chat_domain := mux.Vars(r)["chat-domain"]
	chat, err := h.Services.Repos.Chats.GetChatInfoByDomain(chat_domain)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatByID(w http.ResponseWriter, r *http.Request) {
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		panic(err)
	}
	chat, err := h.Services.Repos.Chats.GetChatInfoByID(chat_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, chat)
}

func (h *Handler) GetChatsByName(w http.ResponseWriter, r *http.Request) {
	chat_name := &models.ChatName{}
	err := json.NewDecoder(r.Body).Decode(&chat_name)
	if err != nil {
		panic(err)
	}
	chat_list, err := h.Services.Repos.Chats.GetListChatsByName(chat_name.Name)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat := &models.CreateChat{}
	err = json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		panic(err)
	}
	chat_id, err := h.Services.Repos.Chats.CreateChat(user_id, chat)
	responder.Respond(w, http.StatusOK, &models.ChatID{ID: chat_id})

}

func (h *Handler) AddUserToChat(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		panic(err)
	}
	err = h.Services.Repos.Chats.AddUserToChat(user_id, chat_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, "")
}

func (h *Handler) GetChatData(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		panic(err)
	}
	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		panic(err)
	}
	chat_data, err := h.Services.Repos.Chats.GetChatDataByID(chat_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, chat_data)
}

func (h *Handler) GetChatMembers(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		panic(err)
	}
	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		panic(err)
	}
	user_list, err := h.Services.Repos.Chats.GetListChatMembers(chat_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) GetChatRooms(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		panic(err)
	}
	if !h.Services.Repos.Chats.UserIsChatMember(user_id, chat_id) {
		panic(err)
	}
	room_list, err := h.Services.Repos.Chats.GetListChatRooms(chat_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, room_list)
}

func (h *Handler) GetUserOwnedChats(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_list, err := h.Services.Repos.Chats.GetChatsOwnedUser(user_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_list, err := h.Services.Repos.Chats.GetChatsInvolvedUser(user_id)
	if err != nil {
		panic(err)
	}
	// json_out, _ := json.MarshalIndent(chat_list, "", "  ")
	// log.Printf("Returning user:\n%s\n", string(json_out))
	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) UpdateChatData(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_id, err := strconv.Atoi(mux.Vars(r)["chat-id"])
	if err != nil {
		panic(err)
	}
	if !h.Services.Repos.Chats.UserIsChatOwner(user_id, chat_id) {
		panic(err)
	}
	chat_data := &models.UpdateChatData{}
	err = json.NewDecoder(r.Body).Decode(&chat_data)
	if err != nil {
		panic(err)
	}
	err = h.Services.Repos.Chats.UpdateChatData(chat_id, chat_data)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, "")
}

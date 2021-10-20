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

func (h *Handler) initUsersRoutes(r *mux.Router) {
	users := r.PathPrefix("/users/").Subrouter()
	{
		// GET
		users.HandleFunc("/d/{user-domain}/", h.GetUserByDomain).Methods(http.MethodGet)
		users.HandleFunc("/{user-id}/", h.GetUserByID).Methods(http.MethodGet)
		users.HandleFunc("/", h.GetUsersByName).Methods(http.MethodGet)

		authenticated := users.PathPrefix("/").Subrouter()
		authenticated.Use(h.checkAuth)
		{
			// GET
			authenticated.HandleFunc("/data/", h.GetUserData).Methods(http.MethodGet)
			authenticated.HandleFunc("/settings/", h.GetUserSettings).Methods(http.MethodGet)
			authenticated.HandleFunc("/chats/owned/", h.GetUserOwnedChats).Methods(http.MethodGet)
			authenticated.HandleFunc("/chats/", h.GetUserChats).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/data/", h.UpdateUserData).Methods(http.MethodPut)
			authenticated.HandleFunc("/settings/", h.UpdateUserSettings).Methods(http.MethodPut)

		}
	}
}

func (h *Handler) GetUserByDomain(w http.ResponseWriter, r *http.Request) {
	user_domain := mux.Vars(r)["user-domain"]
	user, err := h.Services.Repos.Users.GetUserInfoByDomain(user_domain)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		panic(err)
	}
	user, err := h.Services.Repos.Users.GetUserInfoByID(user_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	data, err := h.Services.Repos.Users.GetUserData(user_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, data)
}

func (h *Handler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	settings, err := h.Services.Repos.Users.GetUserSettings(user_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, settings)
}

func (h *Handler) GetUsersByName(w http.ResponseWriter, r *http.Request) {
	name_struct := &models.UserName{}
	err := json.NewDecoder(r.Body).Decode(&name_struct)
	if err != nil {
		panic(err)
	}
	user_list, err := h.Services.Repos.Users.GetListUsersByName(name_struct.Name)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) GetUserOwnedChats(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_list, err := h.Services.Repos.Users.GetListChatsOwnedUser(user_id)
	if err != nil {
		panic(err)
	}
	// json_out, _ := json.MarshalIndent(chat_list, "", "  ")
	// log.Printf("Returning user:\n%s\n", string(json_out))
	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	chat_list, err := h.Services.Repos.Users.GetListChatsUser(user_id)
	if err != nil {
		panic(err)
	}
	// json_out, _ := json.MarshalIndent(chat_list, "", "  ")
	// log.Printf("Returning user:\n%s\n", string(json_out))
	responder.Respond(w, http.StatusOK, chat_list)
}

func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	user_data := &models.UpdateUserData{}
	err = json.NewDecoder(r.Body).Decode(&user_data)
	if err != nil {
		panic(err)
	}
	err = h.Services.Repos.Users.UpdateUserData(user_id, user_data)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, "")
}

func (h *Handler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	user_settings := &models.UpdateUserSettings{}
	err = json.NewDecoder(r.Body).Decode(&user_settings)
	if err != nil {
		panic(err)
	}
	err = h.Services.Repos.Users.UpdateUserSettings(user_id, user_settings)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, "")
}

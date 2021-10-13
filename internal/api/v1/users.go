package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (h *Handler) initUsersRoutes(r *mux.Router) {
	users := r.PathPrefix("/users/").Subrouter()

	users.HandleFunc("/", h.CreateUser).Methods(http.MethodPost)
	users.HandleFunc("/{user-domain}", h.GetUserByDomain).Methods(http.MethodGet)

}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	err = h.Services.Repos.Users.Create(*user)
	if err != nil {
		panic(err)
	}
	user_json, _ := json.MarshalIndent(user, "", "  ")
	log.Printf("%s New user created:\n%s\n", log.Prefix(), string(user_json))
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUserByDomain(w http.ResponseWriter, r *http.Request) {
	// todo: проверка наличия дублирующей записи в бд
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	user_domain := &models.UserDomain{
		Domain: vars["user-domain"],
	}
	user, err := h.Services.Repos.Users.GetByDomain(user_domain)
	if err != nil {
		panic(err)
	}
	user_json, _ := json.MarshalIndent(user, "", "  ")
	log.Printf("%s New user created:\n%s\n", log.Prefix(), string(user_json))
	json.NewEncoder(w).Encode(user)
}

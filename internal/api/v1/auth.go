package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/xlzd/gotp"
)

func (h *Handler) initAuthRoutes(r *mux.Router) {
	auth := r.PathPrefix("/auth/").Subrouter()
	{
		// POST
		auth.HandleFunc("/sign-up/", h.AuthSignUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in/", h.AuthSignIn).Methods(http.MethodPost)
		auth.HandleFunc("/refresh/", h.AuthRefresh).Methods(http.MethodPost)

	}
}

func (h *Handler) AuthSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &models.CreateUser{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	user_id, err := h.Services.Repos.Users.CreateUser(*user)
	if err != nil {
		panic(err)
	}
	uinfo, err := h.Services.Repos.Users.GetUserInfoByID(user_id)
	if err != nil {
		panic(err)
	}
	user_json, _ := json.MarshalIndent(uinfo, "", "  ")
	log.Printf("New user created:\n%s\n", string(user_json))
	responder.Respond(w, http.StatusOK, "")
}

func (h *Handler) AuthSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &models.UserInput{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	user_id, err := h.Services.Repos.Users.GetUserIdByInput(*user)
	if err != nil {
		// todo: uncorrected input
		panic(err)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user-id": user_id,
	}).SignedString(os.Getenv("SECRET_SEGNING_KEY"))
	if err != nil {
		panic(err)
	}
	refresh_token := gotp.RandomSecret(9)
	session := &models.RefreshSession{
		RefreshToken: refresh_token,
		UserAgent:    r.UserAgent(),
		Exp:          int(time.Now().Unix() + int64(time.Hour)), // int(time.Now().Unix() + time.Hour * 24 * 60),
		CreatedAt:    int(time.Now().Unix()),
	}
	err = h.Services.Repos.Users.CreateNewUserRefreshSession(user_id, session)
	if err != nil {
		panic(err)
	}
	token_pair := &models.FreshTokenPair{
		AccessToken:  token,
		RefreshToken: refresh_token,
	}

}

func (h *Handler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
}

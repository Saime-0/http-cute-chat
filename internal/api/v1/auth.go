package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
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
	user_id, err := h.Services.Repos.Users.CreateUser(user)
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
	token_pair, session := GenerateNewSession(user_id, r)
	sessions_count, err := h.Services.Repos.Users.CreateNewUserRefreshSession(user_id, session)
	if err != nil {
		panic(err)
	}
	if sessions_count > 5 {
		err = h.Services.Repos.Users.DeleteOldestSession(user_id)
		if err != nil {
			panic(err)
		}
	}
	responder.Respond(w, http.StatusOK, token_pair)
}

func (h *Handler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	rtoken := &models.TokenForRefreshPair{}
	err := json.NewDecoder(r.Body).Decode(&rtoken)
	if err != nil {
		panic(err)
	}
	session_id, user_id, err := h.Services.Repos.Users.FindSessionByComparedToken(rtoken.RefreshToken)
	if err != nil {
		panic(err)
	}
	token_pair, session := GenerateNewSession(user_id, r)
	err = h.Services.Repos.Users.UpdateRefreshSession(session_id, session)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, token_pair)
}

func GenerateNewSession(user_id int, r *http.Request) (token_pair *models.FreshTokenPair, session *models.RefreshSession) {
	refresh_token := gotp.RandomSecret(16)
	session = &models.RefreshSession{
		RefreshToken: token_pair.RefreshToken,
		UserAgent:    r.UserAgent(),
		Exp:          time.Now().Unix() + int64(time.Hour),
		CreatedAt:    time.Now().Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: session.Exp,
		Subject:   strconv.Itoa(user_id),
	}).SignedString(os.Getenv("SECRET_SEGNING_KEY"))
	if err != nil {
		panic(err)
	}
	token_pair = &models.FreshTokenPair{
		AccessToken:  token,
		RefreshToken: refresh_token,
	}
	return
}

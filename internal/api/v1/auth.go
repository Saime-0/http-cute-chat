package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/saime-0/http-cute-chat/internal/api/rules"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/xlzd/gotp"
)

func (h *Handler) initAuthRoutes(r *mux.Router) {
	auth := r.PathPrefix("/auth").Subrouter()
	{
		// POST
		auth.HandleFunc("/sign-up", h.AuthSignUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.AuthSignIn).Methods(http.MethodPost)
		auth.HandleFunc("/refresh", h.AuthRefresh).Methods(http.MethodPost)

	}
}

func (h *Handler) AuthSignUp(w http.ResponseWriter, r *http.Request) {
	user := &models.CreateUser{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	_, err = h.Services.Repos.Users.CreateUser(user)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) AuthSignIn(w http.ResponseWriter, r *http.Request) {
	user := &models.UserInput{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	user_id, err := h.Services.Repos.Users.GetUserIdByInput(user)
	if err != nil {
		// todo: uncorrected input
		panic(err)
	}
	token_pair, session := GenerateNewSession(user_id, r)
	sessions_count, err := h.Services.Repos.Auth.CreateNewUserRefreshSession(user_id, session)
	if err != nil {
		panic(err)
	}
	if sessions_count > 5 {
		err = h.Services.Repos.Auth.DeleteOldestSession(user_id)
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
	session_id, user_id, err := h.Services.Repos.Auth.FindSessionByComparedToken(rtoken.RefreshToken)
	if err != nil {
		panic(err)
	}
	token_pair, session := GenerateNewSession(user_id, r)
	err = h.Services.Repos.Auth.UpdateRefreshSession(session_id, session)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, token_pair)
}

func GenerateNewSession(user_id int, r *http.Request) (token_pair *models.FreshTokenPair, session *models.RefreshSession) {
	refresh_token := gotp.RandomSecret(rules.RefreshTokenLength)
	session = &models.RefreshSession{
		RefreshToken: refresh_token,
		UserAgent:    r.UserAgent(),
		Exp:          time.Now().Unix() + int64(time.Hour),
		CreatedAt:    time.Now().Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: session.Exp,
		Subject:   strconv.Itoa(user_id),
	}).SignedString([]byte(os.Getenv("SECRET_SIGNING_KEY")))
	if err != nil {
		panic(err)
	}
	token_pair = &models.FreshTokenPair{
		AccessToken:  token,
		RefreshToken: refresh_token,
	}
	return
}

package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/saime-0/http-cute-chat/internal/api/validator"

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

	switch {
	case !validator.ValidateDomain(user.Domain):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidDomain)
		return

	case !validator.ValidateName(user.Name):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidName)
		return

	case !validator.ValidateEmail(user.Email):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidEmail)
		return

	case !validator.ValidatePassword(user.Password):
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidPassword)
		return
	}
	//todo crypt password
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

	userId, err := h.Services.Repos.Users.GetUserIdByInput(user)
	if err != nil {
		// todo: uncorrected input
		panic(err)
	}
	tokenPair, session := GenerateNewSession(userId, r)
	sessionsCount, err := h.Services.Repos.Auth.CreateRefreshSession(userId, session)
	if err != nil {
		panic(err)
	}
	if sessionsCount > 5 {
		err = h.Services.Repos.Auth.DeleteOldestSession(userId)
		if err != nil {
			panic(err)
		}
	}
	responder.Respond(w, http.StatusOK, tokenPair)
}

func (h *Handler) AuthRefresh(w http.ResponseWriter, r *http.Request) {
	rtoken := &models.TokenForRefreshPair{}
	err := json.NewDecoder(r.Body).Decode(&rtoken)
	if err != nil {
		panic(err)
	}
	sessionId, userId, err := h.Services.Repos.Auth.FindSessionByComparedToken(rtoken.RefreshToken)
	if err != nil {
		panic(err)
	}
	tokenPair, session := GenerateNewSession(userId, r)
	err = h.Services.Repos.Auth.UpdateRefreshSession(sessionId, session)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, tokenPair)
}

func GenerateNewSession(userId int, r *http.Request) (tokenPair *models.FreshTokenPair, session *models.RefreshSession) {
	refreshToken := gotp.RandomSecret(rules.RefreshTokenLength)
	session = &models.RefreshSession{
		RefreshToken: refreshToken,
		UserAgent:    r.UserAgent(),
		Exp:          time.Now().Unix() + int64(time.Hour),
		CreatedAt:    time.Now().Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: session.Exp,
		Subject:   strconv.Itoa(userId),
	}).SignedString([]byte(os.Getenv("SECRET_SIGNING_KEY")))
	if err != nil {
		panic(err)
	}
	tokenPair = &models.FreshTokenPair{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	return
}

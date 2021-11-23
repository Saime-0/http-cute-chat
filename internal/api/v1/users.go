package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saime-0/http-cute-chat/internal/api/validator"

	"github.com/saime-0/http-cute-chat/internal/api/rules"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (h *Handler) initUsersRoutes(r *mux.Router) {
	users := r.PathPrefix("/users").Subrouter()
	{
		// ! /users/settings/ equal /users/{user-id}/
		authenticated := users.PathPrefix("/me").Subrouter()
		authenticated.Use(checkAuth)
		{
			// GET
			authenticated.HandleFunc("/data", h.GetUserData).Methods(http.MethodGet)
			authenticated.HandleFunc("/settings", h.GetUserSettings).Methods(http.MethodGet)
			// PUT
			authenticated.HandleFunc("/data", h.UpdateUserData).Methods(http.MethodPut)
			authenticated.HandleFunc("/settings", h.UpdateUserSettings).Methods(http.MethodPut)

		}

		// GET
		users.HandleFunc("/d/{user-domain}", h.GetUserByDomain).Methods(http.MethodGet)
		users.HandleFunc("/{user-id}", h.GetUserByID).Methods(http.MethodGet)
		users.HandleFunc("", h.GetUsersByName).Methods(http.MethodGet)
	}
}

// GetUserByDomain ...
func (h *Handler) GetUserByDomain(w http.ResponseWriter, r *http.Request) {

	userDomain := mux.Vars(r)["user-domain"]
	if !validator.ValidateDomain(userDomain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Users.UserExistsByDomain(userDomain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserNotFound)

		return
	}

	user, err := h.Services.Repos.Users.GetUserByDomain(userDomain)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user)
}

// GetUserByID ...
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(userId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserNotFound)

		return
	}

	user, err := h.Services.Repos.Users.GetUserByID(userId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user)
}

// GetUserData ...
func (h *Handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	data, err := h.Services.Repos.Users.GetUserData(
		r.Context().Value(rules.UserIDFromToken).(int),
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, data)
}

// GetUserSettings ...
func (h *Handler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.Services.Repos.Users.GetUserSettings(
		r.Context().Value(rules.UserIDFromToken).(int),
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, settings)
}

// GetUsersByName ...
func (h *Handler) GetUsersByName(w http.ResponseWriter, r *http.Request) {

	nameFragment := r.URL.Query().Get("name")
	if len(nameFragment) > rules.NameMaxLength || len(nameFragment) == 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	userList, err := h.Services.Repos.Users.GetUsersByNameFragment(nameFragment, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, userList)
}

// UpdateUserData ...
func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	userData := &models.UpdateUserData{}
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	switch {
	case !validator.ValidateDomain(userData.Domain) && len(userData.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidDomain)
		return

	case !validator.ValidateName(userData.Name) && len(userData.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidName)
		return

	case !validator.ValidateEmail(userData.Email) && len(userData.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidEmail)
		return

	case !validator.ValidatePassword(userData.Password) && len(userData.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidPassword)
		return
	}

	if h.Services.Repos.Users.UserExistsByDomain(userData.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	err = h.Services.Repos.Users.UpdateUserData(
		r.Context().Value(rules.UserIDFromToken).(int),
		userData,
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

// UpdateUserSettings ...
func (h *Handler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	userSettings := &models.UpdateUserSettings{}
	err := json.NewDecoder(r.Body).Decode(&userSettings)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if !validator.ValidateAppSettings(userSettings.AppSettings) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Users.UpdateUserSettings(
		r.Context().Value(rules.UserIDFromToken).(int),
		userSettings,
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

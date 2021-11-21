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

	user_domain := mux.Vars(r)["user-domain"]
	if !validator.ValidateDomain(user_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if !h.Services.Repos.Users.UserExistsByDomain(user_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserNotFound)

		return
	}

	user, err := h.Services.Repos.Users.GetUserByDomain(user_domain)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user)
}

// GetUserByID ...
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {

	user_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	if !h.Services.Repos.Users.UserExistsByID(user_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserNotFound)

		return
	}

	user, err := h.Services.Repos.Users.GetUserByID(user_id)
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

	name_fragment := r.URL.Query().Get("name")
	if len(name_fragment) > rules.NameMaxLength || len(name_fragment) == 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}
	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	user_list, err := h.Services.Repos.Users.GetUsersByNameFragment(name_fragment, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user_list)
}

// UpdateUserData ...
func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	user_data := &models.UpdateUserData{}
	err := json.NewDecoder(r.Body).Decode(&user_data)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	switch {
	case !validator.ValidateDomain(user_data.Domain) && len(user_data.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidDomain)
		return

	case !validator.ValidateName(user_data.Name) && len(user_data.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidName)
		return

	case !validator.ValidateEmail(user_data.Email) && len(user_data.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidEmail)
		return

	case !validator.ValidatePassword(user_data.Password) && len(user_data.Domain) != 0:
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidPassword)
		return
	}

	if h.Services.Repos.Users.UserExistsByDomain(user_data.Domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrOccupiedDomain)

		return
	}

	err = h.Services.Repos.Users.UpdateUserData(
		r.Context().Value(rules.UserIDFromToken).(int),
		user_data,
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

// UpdateUserSettings ...
func (h *Handler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	user_settings := &models.UpdateUserSettings{}
	err := json.NewDecoder(r.Body).Decode(&user_settings)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if !validator.ValidateAppSettings(user_settings.AppSettings) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Users.UpdateUserSettings(
		r.Context().Value(rules.UserIDFromToken).(int),
		user_settings,
	)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

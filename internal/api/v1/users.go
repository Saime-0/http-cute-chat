package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *Handler) GetUserByDomain(w http.ResponseWriter, r *http.Request) {
	pl := initPipeline(w, r, h)
	user_domain := pl.parseUserDomainFromRequest()

	if !h.Services.Repos.Users.UserExistsByDomain(user_domain) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserNotFound)

		return
	}

	user, err := h.Services.Repos.Users.GetUserInfoByDomain(user_domain)
	pl.finalInspectionDatabase(err)

	responder.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	pl := initPipeline(w, r, h)

	user_id := pl.parseUserIDFromRequest()

	if !h.Services.Repos.Users.UserExistsByID(user_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrUserNotFound)

		return
	}

	user, err := h.Services.Repos.Users.GetUserInfoByID(user_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	data, err := h.Services.Repos.Users.GetUserData(user_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, data)
}

func (h *Handler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	settings, err := h.Services.Repos.Users.GetUserSettings(user_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, settings)
}

func (h *Handler) GetUsersByName(w http.ResponseWriter, r *http.Request) {
	name_fragment := r.URL.Query().Get("name")
	if len(name_fragment) > rules.NameMaxLength || len(name_fragment) == 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil && r.URL.Query().Get("offset") != "" {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if offset < 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}

	user_list, err := h.Services.Repos.Users.GetUsersByNameFragment(name_fragment, offset)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	user_data := &models.UpdateUserData{}
	err := json.NewDecoder(r.Body).Decode(&user_data)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if (!validateName(user_data.Name) && len(user_data.Name) == 0) && (!validateName(user_data.Domain) && len(user_data.Domain) == 0) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Users.UpdateUserData(user_id, user_data)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	user_settings := &models.UpdateUserSettings{}
	err := json.NewDecoder(r.Body).Decode(&user_settings)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	if !validateAppSettings(user_settings.AppSettings) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	err = h.Services.Repos.Users.UpdateUserSettings(user_id, user_settings)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

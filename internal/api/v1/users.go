package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
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
	user_domain := mux.Vars(r)["user-domain"]
	if !validateDomain(user_domain) {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	user, err := h.Services.Repos.Users.GetUserInfoByDomain(user_domain)
	if err != nil {
		if err == sql.ErrNoRows {
			responder.Error(w, http.StatusNotFound, ErrUserNotFound)
			return
		}
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	user, err := h.Services.Repos.Users.GetUserInfoByID(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			responder.Error(w, http.StatusNotFound, ErrUserNotFound)
			return
		}
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetUserData(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, _ := strconv.Atoi(props["sub"].(string))
	data, err := h.Services.Repos.Users.GetUserData(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			responder.Error(w, http.StatusInternalServerError, ErrDataRetrieved)
			return
		}
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, data)
}

func (h *Handler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, _ := strconv.Atoi(props["sub"].(string))
	settings, err := h.Services.Repos.Users.GetUserSettings(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			responder.Error(w, http.StatusInternalServerError, ErrDataRetrieved)
			return
		}
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, settings)
}

func (h *Handler) GetUsersByName(w http.ResponseWriter, r *http.Request) {
	name_fragment := r.URL.Query().Get("name")
	if len(name_fragment) > NameMaxLength || len(name_fragment) == 0 {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil && r.URL.Query().Get("offset") != "" {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	if offset < 0 {
		responder.Error(w, http.StatusBadRequest, ErrOutOfRange)
		return
	}
	user_list, err := h.Services.Repos.Users.GetUsersByNameFragment(name_fragment, offset)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, _ := strconv.Atoi(props["sub"].(string))
	user_data := &models.UpdateUserData{}
	err := json.NewDecoder(r.Body).Decode(&user_data)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	if (!validateName(user_data.Name) && len(user_data.Name) == 0) && (!validateName(user_data.Domain) && len(user_data.Domain) == 0) {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	err = h.Services.Repos.Users.UpdateUserData(user_id, user_data)
	if err != nil {
		if err == sql.ErrNoRows {
			responder.Error(w, http.StatusInternalServerError, ErrDataRetrieved)
			return
		}
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, _ := strconv.Atoi(props["sub"].(string))
	user_settings := &models.UpdateUserSettings{}
	err := json.NewDecoder(r.Body).Decode(&user_settings)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	if !validateAppSettings(user_settings.AppSettings) {
		responder.Error(w, http.StatusBadRequest, ErrInvalidValue)
		return
	}
	err = h.Services.Repos.Users.UpdateUserSettings(user_id, user_settings)
	if err != nil {
		if err == sql.ErrNoRows {
			responder.Error(w, http.StatusInternalServerError, ErrDataRetrieved)
			return
		}
		responder.Error(w, http.StatusInternalServerError, ErrAccessingDatabase)
		panic(err)
	}
	responder.Respond(w, http.StatusOK, nil)
}

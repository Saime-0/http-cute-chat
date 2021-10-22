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

func (h *Handler) initDialogsRoutes(r *mux.Router) {
	dialogs := r.PathPrefix("/dialogs/").Subrouter()
	{
		//

		authenticated := dialogs.PathPrefix("/").Subrouter()
		authenticated.Use(h.checkAuth)
		{
			// POST
			authenticated.HandleFunc("/{user-id}/messages/", h.SendMessageToUser).Methods(http.MethodPost)
			// GET
			authenticated.HandleFunc("/", h.GetListCompanions).Methods(http.MethodGet)
			authenticated.HandleFunc("/{user-id}/messages/", h.GetListDialogMessages).Methods(http.MethodGet)
			authenticated.HandleFunc("/dialogs/{user-id}/messages/{message-id}/", h.GetDialogMessage).Methods(http.MethodGet)
			// PUT
		}
	}
}

func (h *Handler) SendMessageToUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	target_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		panic(err)
	}
	dialog_id, err := h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, target_id)
	if err != nil {
		panic(err)
	}
	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		panic(err)
	}
	message.Author = user_id
	message_id, err := h.Services.Repos.Dialogs.CreateMessage(dialog_id, message)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, &models.MessageID{ID: message_id})
}

func (h *Handler) GetListCompanions(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	user_list, err := h.Services.Repos.Dialogs.GetCompanions(user_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user_list)

}

func (h *Handler) GetListDialogMessages(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	companion_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		panic(err)
	}
	dialog_id, err := h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, companion_id)
	if err != nil {
		panic(err)
	}
	message_list, err := h.Services.Repos.Dialogs.GetMessages(dialog_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, message_list)
}

func (h *Handler) GetDialogMessage(w http.ResponseWriter, r *http.Request) {
	user_id, err := strconv.Atoi(r.Context().Value("jwt").(jwt.StandardClaims).Subject)
	if err != nil {
		panic(err)
	}
	companion_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		panic(err)
	}
	message_id, err := strconv.Atoi(mux.Vars(r)["message-id"])
	if err != nil {
		panic(err)
	}
	dialog_id, err := h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, companion_id)
	if err != nil {
		panic(err)
	}
	message, err := h.Services.Repos.Dialogs.GetMessageInfo(message_id, dialog_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, message)
}

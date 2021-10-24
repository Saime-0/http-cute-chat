package v1

import (
	"encoding/json"
	"errors"
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
			authenticated.HandleFunc("/{user-id}/messages/", h.GetListDialogMessages).Methods(http.MethodGet)
			authenticated.HandleFunc("/{user-id}/messages/{message-id}/", h.GetDialogMessage).Methods(http.MethodGet)
			authenticated.HandleFunc("/", h.GetListCompanions).Methods(http.MethodGet)
			// PUT
		}
	}
}

func (h *Handler) SendMessageToUser(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
	if err != nil {
		panic(err)
	}
	target_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		panic(err)
	}
	var dialog_id int
	if !h.Services.Repos.Dialogs.DialogIsExistsBetweenUsers(user_id, target_id) {
		dialog_id, err = h.Services.Repos.Dialogs.CreateDialogBetweenUser(user_id, target_id)
		if err != nil {
			panic(err)
		}
	} else {
		dialog_id, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, target_id)
		if err != nil {
			panic(err)
		}
	}

	// todo: create dialog OR create dialog if not exists
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
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
	if err != nil {
		panic(err)
	}
	user_list, err := h.Services.Repos.Dialogs.GetCompanions(user_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, user_list)

}

// ! FIX: при отправке сообщения в диалог, оно попаает не известно куда и при чтении
// ! через GetListDialogMessages возвращается 1 сообщение из комнаты чата
func (h *Handler) GetListDialogMessages(w http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
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
	props, _ := r.Context().Value("jwt").(jwt.MapClaims)
	user_id, err := strconv.Atoi(props["sub"].(string))
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
	var dialog_id int
	if !h.Services.Repos.Dialogs.DialogIsExistsBetweenUsers(user_id, companion_id) {
		panic(errors.New("dialog is undefiend"))
	}
	dialog_id, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, companion_id)
	if err != nil {
		panic(err)
	}
	message, err := h.Services.Repos.Dialogs.GetMessageInfo(message_id, dialog_id)
	if err != nil {
		panic(err)
	}
	responder.Respond(w, http.StatusOK, message)
}

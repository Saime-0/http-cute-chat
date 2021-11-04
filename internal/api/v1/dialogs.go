package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
)

func (h *Handler) initDialogsRoutes(r *mux.Router) {
	authenticated := r.PathPrefix("/dialogs").Subrouter()
	authenticated.Use(checkAuth)
	{
		// POST
		authenticated.HandleFunc("/{user-id}/messages", h.SendMessageToUser).Methods(http.MethodPost)
		// GET
		authenticated.HandleFunc("/{user-id}/messages", h.GetListDialogMessages).Methods(http.MethodGet)
		authenticated.HandleFunc("/{user-id}/messages/{message-id}", h.GetDialogMessage).Methods(http.MethodGet)
		authenticated.HandleFunc("", h.GetListCompanions).Methods(http.MethodGet)
		// PUT
	}
}

func (h *Handler) SendMessageToUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	target_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if target_id == user_id {
		responder.Error(w, http.StatusBadRequest, rules.ErrDialogWithYourself)

		return
	}

	if !h.Services.Repos.Users.UserExistsByID(target_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	var dialog_id int
	if !h.Services.Repos.Dialogs.DialogExistsBetweenUsers(user_id, target_id) {
		dialog_id, err = h.Services.Repos.Dialogs.CreateDialogBetweenUser(user_id, target_id)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	} else {
		dialog_id, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, target_id)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	}

	message := &models.CreateMessage{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrBadRequestBody)

		return
	}

	message.Author = user_id
	message_id, err := h.Services.Repos.Messages.CreateMessageInDialog(dialog_id, message)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.MessageID{ID: message_id})
}

func (h *Handler) GetListCompanions(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	user_list, err := h.Services.Repos.Dialogs.GetCompanions(user_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, user_list)
}

func (h *Handler) GetListDialogMessages(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	target_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if target_id == user_id {
		responder.Error(w, http.StatusBadRequest, rules.ErrDialogWithYourself)

		return
	}

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	if !h.Services.Repos.Users.UserExistsByID(target_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	var dialog_id int
	if !h.Services.Repos.Dialogs.DialogExistsBetweenUsers(user_id, target_id) {
		dialog_id, err = h.Services.Repos.Dialogs.CreateDialogBetweenUser(user_id, target_id)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	} else {
		dialog_id, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, target_id)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	}

	message_list, err := h.Services.Repos.Messages.GetMessagesFromDialog(dialog_id, offset)
	if err != nil {
		panic(err)
	}

	responder.Respond(w, http.StatusOK, message_list)
}

func (h *Handler) GetDialogMessage(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	target_id, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if target_id == user_id {
		responder.Error(w, http.StatusBadRequest, rules.ErrDialogWithYourself)

		return
	}

	if !h.Services.Repos.Users.UserExistsByID(user_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	message_id, err := strconv.Atoi(mux.Vars(r)["message-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	var dialog_id int
	if !h.Services.Repos.Dialogs.DialogExistsBetweenUsers(user_id, target_id) {
		responder.Error(w, http.StatusNotFound, rules.ErrDialogNotFound)

		return
	}

	dialog_id, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(user_id, target_id)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if !h.Services.Repos.Messages.MessageAvailableOnDialog(message_id, dialog_id) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	message, err := h.Services.Repos.Messages.GetMessageFromDialog(message_id, dialog_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, message)
}

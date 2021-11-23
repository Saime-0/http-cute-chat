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
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if targetId == userId {
		responder.Error(w, http.StatusBadRequest, rules.ErrDialogWithYourself)

		return
	}

	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	var dialogId int
	if !h.Services.Repos.Dialogs.DialogExistsBetweenUsers(userId, targetId) {
		dialogId, err = h.Services.Repos.Dialogs.CreateDialogBetweenUser(userId, targetId)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	} else {
		dialogId, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(userId, targetId)
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

	message.Author = userId
	messageId, err := h.Services.Repos.Messages.CreateMessageInDialog(dialogId, message)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, &models.MessageID{ID: messageId})
}

func (h *Handler) GetListCompanions(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	userList, err := h.Services.Repos.Dialogs.GetCompanions(userId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, userList)
}

func (h *Handler) GetListDialogMessages(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if targetId == userId {
		responder.Error(w, http.StatusBadRequest, rules.ErrDialogWithYourself)

		return
	}

	offset, ok := parseOffsetFromQuery(w, r)
	if !ok {

		return
	}

	if !h.Services.Repos.Users.UserExistsByID(targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	var dialogId int
	if !h.Services.Repos.Dialogs.DialogExistsBetweenUsers(userId, targetId) {
		dialogId, err = h.Services.Repos.Dialogs.CreateDialogBetweenUser(userId, targetId)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	} else {
		dialogId, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(userId, targetId)
		if err != nil {
			responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
			panic(err)
		}
	}

	messageList, err := h.Services.Repos.Messages.GetMessagesFromDialog(dialogId, offset)
	if err != nil {
		panic(err)
	}

	responder.Respond(w, http.StatusOK, messageList)
}

func (h *Handler) GetDialogMessage(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	targetId, err := strconv.Atoi(mux.Vars(r)["user-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if targetId == userId {
		responder.Error(w, http.StatusBadRequest, rules.ErrDialogWithYourself)

		return
	}

	if !h.Services.Repos.Users.UserExistsByID(userId) {
		responder.Error(w, http.StatusNotFound, rules.ErrUserNotFound)

		return
	}

	messageId, err := strconv.Atoi(mux.Vars(r)["message-id"])
	if err != nil {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	var dialogId int
	if !h.Services.Repos.Dialogs.DialogExistsBetweenUsers(userId, targetId) {
		responder.Error(w, http.StatusNotFound, rules.ErrDialogNotFound)

		return
	}

	dialogId, err = h.Services.Repos.Dialogs.GetDialogIDBetweenUsers(userId, targetId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}

	if !h.Services.Repos.Messages.MessageAvailableOnDialog(messageId, dialogId) {
		responder.Error(w, http.StatusBadRequest, rules.ErrNoAccess)

		return
	}

	message, err := h.Services.Repos.Messages.GetMessageFromDialog(messageId, dialogId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, message)
}

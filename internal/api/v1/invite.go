package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

func (h *Handler) initInviteRoutes(r *mux.Router) {
	invites := r.PathPrefix("/invites").Subrouter()
	{
		authenticated := invites.PathPrefix("").Subrouter()
		authenticated.Use(checkAuth)
		{
			// POST
			authenticated.HandleFunc("{invite-code}", h.JoinToChatByCode).Methods(http.MethodPost)
		}
		// GET
		invites.HandleFunc("{invite-code}", h.GetChatByCode).Methods(http.MethodGet)
	}
}

func (h *Handler) JoinToChatByCode(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(rules.UserIDFromToken).(int)

	linkCode := mux.Vars(r)["invite-code"]
	if !h.Services.Repos.Chats.InviteLinkIsRelevant(linkCode) {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	l, err := h.Services.Repos.Chats.FindInviteLinkByCode(linkCode)
	if err != nil {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	if h.Services.Repos.Chats.UserIsChatMember(userId, l.ChatID) {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	countChats, err := h.Services.Repos.Chats.GetCountUserChats(userId)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if countChats >= rules.MaxUserChats {
		responder.Error(w, http.StatusBadRequest, rules.ErrLimitHasBeenReached)

		return
	}

	countMembers, err := h.Services.Repos.Chats.GetCountChatMembers(l.ChatID)
	if err != nil {
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)

		panic(err)
	}
	if countMembers >= rules.MaxMembersOnChat {
		responder.Error(w, http.StatusBadRequest, rules.ErrMembersLimitHasBeenReached)

		return
	}

	_, err = h.Services.Repos.Chats.AddUserByCode(linkCode, userId)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatByCode(w http.ResponseWriter, r *http.Request) {
	linkCode := mux.Vars(r)["invite-code"]
	if !h.Services.Repos.Chats.InviteLinkIsRelevant(linkCode) {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	l, err := h.Services.Repos.Chats.FindInviteLinkByCode(linkCode)
	if err != nil {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	chat, err := h.Services.Repos.Chats.GetChatByID(l.ChatID)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat)
}

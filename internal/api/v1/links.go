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
			authenticated.HandleFunc("{link-code}", h.JoinToChatByCode).Methods(http.MethodPost)
		}
		// GET
		invites.HandleFunc("{link-code}", h.GetChatByCode).Methods(http.MethodGet)
	}
}

func (h *Handler) JoinToChatByCode(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(rules.UserIDFromToken).(int)

	link_code := mux.Vars(r)["link-code"]
	if !h.Services.Repos.Chats.InviteLinkIsRelevant(link_code) {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	_, err := h.Services.Repos.Chats.AddUserByCode(link_code, user_id)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, nil)
}

func (h *Handler) GetChatByCode(w http.ResponseWriter, r *http.Request) {
	link_code := mux.Vars(r)["link-code"]
	if !h.Services.Repos.Chats.InviteLinkIsRelevant(link_code) {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}

	l, err := h.Services.Repos.Chats.FindInviteLinkByCode(link_code)
	if err != nil {
		responder.Error(w, http.StatusNotFound, rules.ErrInviteLinkNotFound)

		return
	}
	chat, err := h.Services.Repos.Chats.GetChatByID(l.ChatID)
	finalInspectionDatabase(w, err)

	responder.Respond(w, http.StatusOK, chat)
}

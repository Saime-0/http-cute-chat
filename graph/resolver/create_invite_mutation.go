package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) CreateInvite(ctx context.Context, input model.CreateInviteInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > CreateInvite [cid:", input.ChatID, "]")
	defer node.Kill()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(input.ChatID) ||
		node.IsMember(clientID, input.ChatID) ||
		node.CanCreateInvite(clientID, input.ChatID) ||
		node.InvitesLimit(input.ChatID) ||
		node.ValidInviteInput(input) {
		return node.Err, nil
	}

	eventReadyInvite, err := r.Services.Repos.Chats.CreateInvite(&input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать инвайт"), nil
	}
	go r.Subix.NotifyChatMembers(
		input.ChatID,
		eventReadyInvite,
	)

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); eventReadyInvite.ExpiresAt != nil && ok && *eventReadyInvite.ExpiresAt < runAt.(int64) {
		r.CreateScheduledInvite(input.ChatID, eventReadyInvite.Code, eventReadyInvite.ExpiresAt)
	}

	return resp.Success("инвайт создан"), nil
}

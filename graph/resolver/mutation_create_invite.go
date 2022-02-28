package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateInvite(ctx context.Context, input model.CreateInviteInput) (model.CreateInviteResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("CreateInvite", &bson.M{
		"input": input,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(input.ChatID) ||
		node.IsMember(clientID, input.ChatID) ||
		node.CanCreateInvite(clientID, input.ChatID) ||
		node.InvitesLimit(input.ChatID) ||
		node.ValidInviteInput(input) {
		return node.GetError(), nil
	}

	eventReadyInvite, err := r.Services.Repos.Chats.CreateInvite(&input)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}
	go r.Subix.NotifyChatMembers(
		input.ChatID,
		eventReadyInvite,
	)

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); eventReadyInvite.ExpiresAt != nil && ok && *eventReadyInvite.ExpiresAt < runAt.(int64) {
		err := r.CreateScheduledInvite(input.ChatID, eventReadyInvite.Code, eventReadyInvite.ExpiresAt)
		if err != nil {
			node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		}
	}

	return &model.CreatedInvite{InviteCode: eventReadyInvite.Code}, nil
}

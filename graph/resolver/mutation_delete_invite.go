package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) DeleteInvite(ctx context.Context, chatID int, code string) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("DeleteInvite", &bson.M{
		"chatID": chatID,
		"code":   code,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(chatID) ||
		node.IsMember(clientID, chatID) ||
		node.CanCreateInvite(clientID, chatID) ||
		node.HasInvite(chatID, code) ||
		node.InviteIsRelevant(code) {
		return node.GetError(), nil
	}

	eventReadyInvite, err := r.Services.Repos.Chats.DeleteInvite(code)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось удалить инвайт"), nil
	}

	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyInvite,
	)
	r.ForceDropScheduledInvite(eventReadyInvite.Code) // да, я знаю что это криндж..

	return resp.Success("инвайт успешно удален"), nil
}

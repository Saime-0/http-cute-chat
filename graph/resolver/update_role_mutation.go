package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) UpdateRole(ctx context.Context, roleID int, input model.UpdateRoleInput) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("UpdateRole", &bson.M{
		"roleID": roleID,
		"input":  input,
	})
	defer node.MethodTiming()

	var (
		chatID   int
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
	)

	if input.Name != nil && node.ValidName(*input.Name) ||
		node.GetChatIDByRole(roleID, &chatID) ||
		node.CanUpdateRole(clientID, chatID) {
		return node.GetError(), nil
	}

	eventReadyRole, err := r.Services.Repos.Chats.UpdateRole(roleID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные"), nil
	}

	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyRole,
	)

	return resp.Success("данные успешно обновлены"), nil
}

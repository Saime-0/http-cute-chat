package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateRole(ctx context.Context, input model.CreateRoleInput) (model.CreateRoleResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("CreateRole", &bson.M{
		"input": input,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(input.ChatID) ||
		node.IsMember(clientID, input.ChatID) ||
		node.CanCreateRole(clientID, input.ChatID) ||
		node.RolesLimit(input.ChatID) {
		return node.GetError(), nil
	}

	eventReadyRole, err := r.Services.Repos.Chats.CreateRoleInChat(&input)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	go r.Subix.NotifyChatMembers(
		input.ChatID,
		eventReadyRole,
	)

	return &model.CreatedRole{RoleID: eventReadyRole.ID}, nil
}

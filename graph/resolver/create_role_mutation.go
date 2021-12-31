package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
)

func (r *mutationResolver) CreateRole(ctx context.Context, input model.CreateRoleInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > CreateRole [cid:", input.ChatID, "]")
	defer node.Kill()

	clientID := ctx.Value(rules.UserIDFromToken).(int)

	if node.ChatExists(input.ChatID) ||
		node.IsMember(clientID, input.ChatID) ||
		node.CanCreateRole(clientID, input.ChatID) ||
		node.RolesLimit(input.ChatID) {
		return node.Err, nil
	}

	err := r.Services.Repos.Chats.CreateRoleInChat(&input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать роль"), nil
	}

	return resp.Success("роль создана"), nil
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) DeleteRole(ctx context.Context, roleID int) (model.MutationResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   int
	)

	if node.GetChatIDByRole(roleID, &chatID) ||
		node.CanCreateRole(clientID, chatID) {
		return node.Err, nil
	}

	eventReadyRole, err := r.Services.Repos.Chats.DeleteRole(roleID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось удалить роль"), nil
	}

	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyRole,
	)

	return resp.Success("роль удалена"), nil
}

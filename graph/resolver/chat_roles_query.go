package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *queryResolver) ChatRoles(ctx context.Context, chatID int) (model.ChatRolesResult, error) {
	node := r.Piper.CreateNode("queryResolver > ChatRoles [cid:", chatID, "]")
	defer node.Kill()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.ChatExists(chatID) ||
		node.IsMember(clientID, chatID) {
		return node.Err, nil
	}

	roles, err := r.Services.Repos.Chats.Roles(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось получить список ролей в комнате"), nil
	}

	return roles, nil
}

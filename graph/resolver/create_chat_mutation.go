package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *mutationResolver) CreateChat(ctx context.Context, input model.CreateChatInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > CreateChat [_]")
	defer node.Kill()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.OwnedLimit(clientID) ||
		node.ChatsLimit(clientID) ||
		node.ValidDomain(input.Domain) ||
		node.ValidName(input.Name) ||
		node.DomainIsFree(input.Domain) {
		return node.Err, nil
	}

	chatID, err := r.Services.Repos.Chats.CreateChat(clientID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось создать чат"), nil
	}

	_, err = r.Services.Repos.Chats.AddUserToChat(clientID, chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось присоединиться к чату"), nil
	}

	return resp.Success("чат успешно создан"), nil
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/internal/utils"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/resp"
)

func (r *mutationResolver) UpdateChat(ctx context.Context, chatID int, input model.UpdateChatInput) (model.MutationResult, error) {
	node := r.Piper.CreateNode("mutationResolver > UpdateChat [cid:", chatID, "]")
	defer node.Kill()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.CanUpdateChat(clientID, chatID) ||
		input.Name != nil && node.ValidName(*input.Name) ||
		input.Domain != nil && node.ValidDomain(*input.Domain) {
		return node.Err, nil
	}

	eventReadyChat, err := r.Services.Repos.Chats.UpdateChat(chatID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные чата"), nil
	}
	go r.Services.Subix.NotifyChatMembers(
		chatID,
		eventReadyChat,
	)
	return resp.Success("данные чата обновлены"), nil
}

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

func (r *mutationResolver) UpdateChat(ctx context.Context, chatID int, input model.UpdateChatInput) (model.MutationResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("UpdateChat", &bson.M{
		"chatID": chatID,
		"input":  input,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.CanUpdateChat(clientID, chatID) ||
		input.Name != nil && node.ValidName(*input.Name) ||
		input.Domain != nil && node.ValidDomain(*input.Domain) {
		return node.GetError(), nil
	}

	eventReadyChat, err := r.Services.Repos.Chats.UpdateChat(chatID, &input)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить данные чата"), nil
	}
	go r.Subix.NotifyChatMembers(
		chatID,
		eventReadyChat,
	)
	return resp.Success("данные чата обновлены"), nil
}

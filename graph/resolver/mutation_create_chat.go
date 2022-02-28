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

func (r *mutationResolver) CreateChat(ctx context.Context, input model.CreateChatInput) (model.CreateChatResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("CreateChat", &bson.M{
		"input": input,
	})
	defer node.MethodTiming()

	clientID := utils.GetAuthDataFromCtx(ctx).UserID

	if node.OwnedLimit(clientID) ||
		node.ChatsLimit(clientID) ||
		node.ValidDomain(input.Domain) ||
		node.ValidName(input.Name) ||
		node.DomainIsFree(input.Domain) {
		return node.GetError(), nil
	}

	chatID, err := r.Services.Repos.Chats.CreateChat(clientID, &input)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось создать чат"), nil
	}

	_, err = r.Services.Repos.Chats.AddUserToChat(clientID, chatID)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "не удалось присоединиться к чату"), nil
	}

	return &model.CreatedChat{ChatID: chatID}, nil
}

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *queryResolver) Messages(ctx context.Context, find model.FindMessages, params *model.Params) (model.MessagesResult, error) {
	node := r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	var (
		clientID = utils.GetAuthDataFromCtx(ctx).UserID
		chatID   = find.ChatID
		holder   models.AllowHolder
		messages *model.Messages
	)

	if node.ValidParams(&params) ||
		node.ValidID(chatID) ||
		node.IsMember(clientID, chatID) ||
		find.RoomID != nil && node.ValidID(*find.RoomID) ||
		find.UserID != nil && node.ValidID(*find.UserID) ||
		node.GetAllowHolder(clientID, chatID, &holder) { // todo bodyfragment valid
		return node.Err, nil
	}

	messages = r.Services.Repos.Chats.FindMessages(&find, params, &holder)
	return messages, nil
}

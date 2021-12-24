package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/tlog"
)

func (r *queryResolver) Messages(ctx context.Context, find model.FindMessages, params model.Params) (model.MessagesResult, error) {
	tl := tlog.Start("queryResolver > Messages [cid:", find.ChatID, "]")
	defer tl.Fine()

	clientID := ctx.Value(rules.UserIDFromToken).(int)

	node := r.Piper.CreateNode()
	defer node.Kill()

	var (
		chatID   = find.ChatID
		holder   models.AllowHolder
		messages *model.Messages
	)

	if node.ValidParams(&params) ||
		node.ValidID(chatID) ||
		node.IsMember(clientID, chatID) ||
		find.RoomID != nil && node.ValidID(*find.RoomID) ||
		find.AuthorID != nil && node.ValidID(*find.AuthorID) ||
		node.GetAllowHolder(clientID, chatID, &holder) { // todo bodyfragment valid
		return node.Err, nil
	}

	messages = r.Services.Repos.Chats.FindMessages(&find, &params, &holder)
	return messages, nil
}

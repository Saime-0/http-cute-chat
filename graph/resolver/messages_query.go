package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Messages(ctx context.Context, find model.FindMessages, params *model.Params) (model.MessagesResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	var (
		chatID   int
		holder   models.AllowHolder
		messages *model.Messages
	)
	if pl.ValidID(find.ChatID) ||
		pl.IsMember(clientID, find.ChatID) ||
		find.RoomID != nil && pl.ValidID(*find.RoomID) ||
		find.AuthorID != nil && pl.ValidID(*find.AuthorID) ||
		pl.GetAllowHolder(clientID, chatID, &holder) {
		return pl.Err, nil
	}

	messages = r.Services.Repos.Chats.FindMessages(&find, &holder)
	return messages, nil
}

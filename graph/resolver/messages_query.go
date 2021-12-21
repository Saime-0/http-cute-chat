package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/saime-0/http-cute-chat/internal/tlog"
	"strconv"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Messages(ctx context.Context, find model.FindMessages, params model.Params) (model.MessagesResult, error) {
	tl := tlog.Start("queryResolver > Messages [cid:" + strconv.Itoa(find.ChatID) + "]")
	defer tl.Fine()
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var (
		chatID   = find.ChatID
		holder   models.AllowHolder
		messages *model.Messages
	)

	if pl.ValidParams(&params) ||
		pl.ValidID(chatID) ||
		pl.IsMember(clientID, chatID) ||
		find.RoomID != nil && pl.ValidID(*find.RoomID) ||
		find.AuthorID != nil && pl.ValidID(*find.AuthorID) ||
		pl.GetAllowHolder(clientID, chatID, &holder) { // todo bodyfragment valid
		return pl.Err, nil
	}

	messages = r.Services.Repos.Chats.FindMessages(&find, &params, &holder)
	return messages, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func dbg(str string) bool {
	fmt.Println(str)
	return true
}

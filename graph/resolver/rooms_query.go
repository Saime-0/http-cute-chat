package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
	"github.com/saime-0/http-cute-chat/internal/tlog"
)

func (r *queryResolver) Rooms(ctx context.Context, find model.FindRooms, params model.Params) (model.RoomsResult, error) {
	tl := tlog.Start("queryResolver > Rooms [cid:" + strconv.Itoa(find.ChatID) + "]")
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(r.Services.Repos)
	var (
		chatID = find.ChatID
		rooms  *model.Rooms
	)

	if pl.ValidParams(&params) ||
		pl.ValidID(chatID) ||
		pl.IsMember(clientID, chatID) ||
		find.RoomID != nil && pl.ValidID(*find.RoomID) ||
		find.ParentID != nil && pl.ValidID(*find.ParentID) {
		tl.FineWithReason(pl.Err.Error)
		return pl.Err, nil
	}

	rooms = r.Services.Repos.Rooms.FindRooms(&find, &params)
	tl.Fine()
	return rooms, nil
}

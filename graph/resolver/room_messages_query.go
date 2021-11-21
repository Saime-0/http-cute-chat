package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *queryResolver) RoomMessages(ctx context.Context, roomID *int64, filter *model.MessageFilter, params *model.Params) (model.RoomMessagesResult, error) {
	panic(fmt.Errorf("not implemented"))
}

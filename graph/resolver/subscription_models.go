package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
)

func (r *updateChatResolver) CountMembers(ctx context.Context, obj *model.UpdateChat) (model.CountMembersResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// UpdateChat returns generated.UpdateChatResolver implementation.
func (r *Resolver) UpdateChat() generated.UpdateChatResolver { return &updateChatResolver{r} }

type updateChatResolver struct{ *Resolver }

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Chat(ctx context.Context, input model.FindByDomainOrID) (model.ChatResult, error) {
	//clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if input.ID != nil && pl.ChatExists(*input.ID) ||
		input.Domain != nil && pl.ChatExistsByDomain(*input.Domain) &&
			pl.GetIDByDomain(*input.Domain, input.ID) {
		return pl.Err, nil
	}
	chatData, err := r.Services.Repos.Chats.GetChatByID(*input.ID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return model.Chat{
		Unit: &model.Unit{
			ID:     chatData.ID,
			Domain: chatData.Domain,
			Name:   chatData.Name,
			Type:   model.UnitTypeChat,
		},
		Owner:        nil, // forced
		Rooms:        nil, // forced
		Private:      chatData.Private,
		CountMembers: 0,   // forced
		Members:      nil, // forced
		Roles:        nil, // forced
		Invites:      nil, // forced
		Banlist:      nil, // forced
		MeRestricts:  nil, // forced
	}, nil
}

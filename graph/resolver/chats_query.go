package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) Chats(ctx context.Context, nameFragment string, params *model.Params) (model.ChatsResult, error) {
	//clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.ValidParams(params) ||
		pl.ValidNameFragment(nameFragment) {
		return pl.Err, nil
	}

	chats, err := r.Services.Repos.Chats.GetChatsByNameFragment(nameFragment, *params.Limit, *params.Offset)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	m := model.Chats{
		Chats: []*model.Chat{},
	}
	for _, chat := range chats {
		m.Chats = append(m.Chats, &model.Chat{
			Unit: &model.Unit{
				ID:     chat.Unit.ID,
				Domain: chat.Unit.Domain,
				Name:   chat.Unit.Name,
				Type:   model.UnitType(chat.Unit.Type),
			},
			Owner:        nil, // forced
			Rooms:        nil, // forced
			Private:      chat.Private,
			CountMembers: nil, // forced
			Members:      nil, // forced
			Roles:        nil, // forced
			Invites:      nil, // forced
			Banlist:      nil, // forced
			Me:           nil, // forced
		})
	}

	return m, nil
}

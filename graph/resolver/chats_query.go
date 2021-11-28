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

	m := model.ChatArray{
		Chats: []*model.Chat{},
	}
	for _, chat := range chats {
		m.Chats = append(m.Chats, &model.Chat{
			Unit: &model.Unit{
				ID:     chat.ID,
				Domain: chat.Domain,
				Name:   chat.Name,
				Type:   model.UnitTypeChat,
			},
			Owner:        nil, // forced
			Rooms:        nil, // forced
			Private:      chat.Private,
			CountMembers: 0,   // forced
			Members:      nil, // forced
			Roles:        nil, // forced
			Invites:      nil, // forced
			Banlist:      nil, // forced
			MeRestricts:  nil, // forced
		})
	}

	return m, nil
}

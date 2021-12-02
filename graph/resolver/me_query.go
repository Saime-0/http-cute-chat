package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

func (r *meResolver) Chats(ctx context.Context, obj *model.Me) ([]*model.Chat, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	_chats, err := r.Services.Repos.Users.Chats(clientID)
	if err != nil {
		return nil, nil // resp.Error
	}
	var chats []*model.Chat
	for _, chat := range _chats {
		chats = append(chats, &model.Chat{
			Unit: &model.Unit{
				ID:     chat.Unit.ID,
				Domain: chat.Unit.Domain,
				Name:   chat.Unit.Name,
				Type:   model.UnitType(chat.Unit.Type),
			},
			Private: chat.Private,
		})
	}
	return chats, nil
}

func (r *meResolver) OwnedChats(ctx context.Context, obj *model.Me) ([]*model.Chat, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	_chats, err := r.Services.Repos.Users.OwnedChats(clientID)
	if err != nil {
		return nil, nil // resp.Error
	}
	var chats []*model.Chat
	for _, chat := range _chats {
		chats = append(chats, &model.Chat{
			Unit: &model.Unit{
				ID:     chat.Unit.ID,
				Domain: chat.Unit.Domain,
				Name:   chat.Unit.Name,
				Type:   model.UnitType(chat.Unit.Type),
			},
			Private: chat.Private,
		})
	}
	return chats, nil
}

func (r *queryResolver) Me(ctx context.Context) (model.MeResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	me, err := r.Services.Repos.Users.Me(clientID)
	if err != nil {
		println(err.Error())
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	return model.Me{
		User: &model.User{
			Unit: &model.Unit{
				ID:     me.Unit.ID,
				Domain: me.Unit.Domain,
				Name:   me.Unit.Name,
				Type:   model.UnitType(me.Unit.Type), // 0_o
			},
		},
		Data: &model.UserData{
			Email:    me.User.Email,
			Password: me.User.Password,
		},
		Chats:      nil, // forced
		OwnedChats: nil, // forced
	}, nil
}

// Me returns generated.MeResolver implementation.
func (r *Resolver) Me() generated.MeResolver { return &meResolver{r} }

type meResolver struct{ *Resolver }

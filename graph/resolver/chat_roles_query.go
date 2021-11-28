package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/piping"
)

func (r *queryResolver) ChatRoles(ctx context.Context, chatID int) (model.ChatRolesResult, error) {
	clientID := ctx.Value(rules.UserIDFromToken).(int)
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	if pl.ChatExists(chatID) ||
		pl.IsMember(clientID, chatID) {
		return pl.Err, nil
	}

	roles, err := r.Services.Repos.Chats.ChatRoles(chatID)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}

	m := model.RoleArray{
		Roles: []*model.Role{},
	}
	for _, role := range roles {
		m.Roles = append(m.Roles, &model.Role{
			ID:    role.ID,
			Users: nil, // forced
			Name:  role.Name,
			Color: role.Color,
		})
	}

	return m, nil
}

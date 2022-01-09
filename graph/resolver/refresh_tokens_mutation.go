package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/pkg/kit"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
)

func (r *mutationResolver) RefreshTokens(ctx context.Context, refreshToken string) (model.RefreshTokensResult, error) {
	node := r.Piper.CreateNode("mutationResolver > RefreshTokens [<token>]")
	defer node.Kill()

	_, clientID, err := r.Services.Repos.Auth.FindSessionByComparedToken(refreshToken)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обрабработать токен"), nil
	}

	var (
		session *models.RefreshSession
	)
	newRefreshToken := kit.CryptoSecret(rules.RefreshTokenBytesLength)
	session = &models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserAgent:    ctx.Value(rules.UserAgentFromHeaders).(string),
		Lifetime:     rules.RefreshTokenLiftime,
	}

	expiresAt, err := r.Services.Repos.Auth.CreateRefreshSession(clientID, session, true)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить сессию"), nil
	}

	token, err := utils.GenerateToken(
		&utils.TokenData{
			UserID:    clientID,
			ExpiresAt: expiresAt,
		},
		r.Config.SecretKey,
	)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при обработке токена"), nil
	}

	return model.TokenPair{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}, nil
}

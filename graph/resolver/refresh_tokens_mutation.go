package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/xlzd/gotp"
)

func (r *mutationResolver) RefreshTokens(ctx context.Context, refreshToken string) (model.RefreshTokensResult, error) {
	_, clientID, err := r.Services.Repos.Auth.FindSessionByComparedToken(refreshToken)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обрабработать токен"), nil
	}

	var (
		session *models.RefreshSession
	)
	newRefreshToken := gotp.RandomSecret(rules.RefreshTokenLength)
	session = &models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserAgent:    ctx.Value(rules.UserAgentFromHeaders).(string),
		Lifetime:     rules.RefreshTokenLiftime,
	}
	// todo UpdaateRefreshSession
	// fixme удаление не тех записией
	expiresAt, err := r.Services.Repos.Auth.CreateRefreshSession(clientID, session, true)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить сессию"), nil
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: expiresAt,
		Subject:   strconv.Itoa(clientID),
	}).SignedString([]byte(os.Getenv("SECRET_SIGNING_KEY")))
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при обработке токена"), nil
	}

	return model.TokenPair{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}, nil
}

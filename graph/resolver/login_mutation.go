package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (model.LoginResult, error) {
	node := r.Piper.CreateNode("mutationResolver > Login [_]")
	defer node.Kill()

	var clientID int

	if node.UserExistsByInput(input) ||
		node.GetUserIDByInput(input, &clientID) {
		return node.Err, nil
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
		return resp.Error(resp.ErrInternalServerError, "неудачная попытка создать сессию пользователя"), nil
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

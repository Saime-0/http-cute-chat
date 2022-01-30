package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (r *mutationResolver) RefreshTokens(ctx context.Context, sessionKey *string, refreshToken string) (model.RefreshTokensResult, error) {
	node := r.Piper.CreateNode("mutationResolver > RefreshTokens [<token>]")
	defer node.Kill()

	sessionID, clientID, err := r.Services.Repos.Auth.FindSessionByComparedToken(refreshToken)
	if err != nil {
		println("RefreshTokens:", err.Error()) // debug
		return resp.Error(resp.ErrInternalServerError, "не удалось обрабработать токен"), nil
	}

	var (
		session *models.RefreshSession
	)
	newRefreshToken := kit.RandomSecret(rules.RefreshTokenLength)
	sessionExpAt := kit.After(rules.RefreshTokenLiftime)
	session = &models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserAgent:    ctx.Value(res.CtxUserAgent).(string),
		ExpAt:        sessionExpAt,
	}

	err = r.Services.Repos.Auth.UpdateRefreshSession(sessionID, session)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "не удалось обновить сессию"), nil
	}

	tokenExpiresAt := kit.After(rules.AccessTokenLiftime)
	token, err := utils.GenerateToken(
		&utils.TokenData{
			UserID:    clientID,
			ExpiresAt: tokenExpiresAt,
		},
		r.Config.SecretKey,
	)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при обработке токена"), nil
	}

	if sessionKey != nil {
		err = r.Services.Subix.ExtendClientSession(*sessionKey, tokenExpiresAt)
		if err != nil {
			println("RefreshTokens:", err) // debug
		}
	}

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); ok && sessionExpAt < runAt.(int64) {
		_, err = r.Services.Scheduler.AddTask(
			func() {
				r.Services.Repos.Users.DeleteRefreshSession(sessionID)
			},
			sessionExpAt,
		)
		if err != nil {
			panic(err)
		}
	}

	return model.TokenPair{
		AccessToken:  token,
		RefreshToken: newRefreshToken,
	}, nil
}

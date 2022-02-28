package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/kit"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) RefreshTokens(ctx context.Context, sessionKey *string, refreshToken string) (model.RefreshTokensResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("RefreshTokens", &bson.M{
		"sessionKey":   sessionKey,
		"refreshToken": refreshToken,
	})
	defer node.MethodTiming()

	sessionID, clientID, err := r.Services.Repos.Auth.FindSessionByComparedToken(refreshToken)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	var (
		session *models.RefreshSession
	)
	newRefreshToken := kit.RandomSecret(rules.RefreshTokenLength)
	sessionExpAt := kit.After(*r.Config.RefreshTokenLiftime)
	session = &models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserAgent:    ctx.Value(res.CtxUserAgent).(string),
		ExpAt:        sessionExpAt,
	}

	err = r.Services.Repos.Auth.UpdateRefreshSession(sessionID, session)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	tokenExpiresAt := kit.After(*r.Config.AccessTokenLiftime)
	token, err := utils.GenerateToken(
		&utils.TokenData{
			UserID:    clientID,
			ExpiresAt: tokenExpiresAt,
		},
		r.Config.SecretSigningKey,
	)
	if err != nil {
		node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		return resp.Error(resp.ErrInternalServerError, "произошла ошибка во время обработки данных"), nil
	}

	if sessionKey != nil {
		err = r.Subix.ExtendClientSession(*sessionKey, tokenExpiresAt)
		if err != nil {
			node.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
		}
	}

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); ok && sessionExpAt < runAt.(int64) {
		_, err = r.Services.Scheduler.AddTask(
			func() {
				err := r.Services.Repos.Users.DeleteRefreshSession(sessionID)
				if err != nil {
					r.Healer.Alert(cerrors.Wrap(err, utils.GetCallerPos()))
				}
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

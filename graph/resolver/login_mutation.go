package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/resp"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/kit"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (model.LoginResult, error) {
	node := *r.Piper.NodeFromContext(ctx)
	defer r.Piper.DeleteNode(*node.ID)

	node.SwitchMethod("Login")
	defer node.MethodTiming()

	var (
		clientID   int
		requisites = &models.LoginRequisites{
			Email: input.Email,
			HashedPasswd: func() string {
				hpasswd, err := utils.HashPassword(input.Password, r.Config.PasswordSalt)
				if err != nil {
					panic(err)
				}
				return hpasswd
			}(),
		}
	)

	if node.UserExistsByRequisites(requisites) ||
		node.GetUserIDByRequisites(requisites, &clientID) {
		return node.Err, nil
	}

	var (
		session *models.RefreshSession
	)
	newRefreshToken := kit.RandomSecret(rules.RefreshTokenLength)
	expAt := kit.After(rules.RefreshTokenLiftime)
	session = &models.RefreshSession{
		RefreshToken: newRefreshToken,
		UserAgent:    ctx.Value(res.CtxUserAgent).(string),
		ExpAt:        expAt,
	}
	sessionID, err := r.Services.Repos.Auth.CreateRefreshSession(clientID, session, true)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "неудачная попытка создать сессию пользователя"), nil
	}

	token, err := utils.GenerateToken(
		&utils.TokenData{
			UserID:    clientID,
			ExpiresAt: time.Now().Unix() + rules.AccessTokenLiftime,
		},
		r.Config.SecretKey,
	)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "ошибка при обработке токена"), nil
	}

	if runAt, ok := r.Services.Cache.Get(res.CacheNextRunRegularScheduleAt); ok && expAt < runAt.(int64) {
		_, err = r.Services.Scheduler.AddTask(
			func() {
				r.Services.Repos.Users.DeleteRefreshSession(sessionID)
			},
			expAt,
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

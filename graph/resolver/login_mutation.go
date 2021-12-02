package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/api/resp"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"github.com/saime-0/http-cute-chat/internal/models"
	"github.com/saime-0/http-cute-chat/internal/piping"
	"github.com/xlzd/gotp"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (model.LoginResult, error) {
	pl := piping.NewPipeline(ctx, r.Services.Repos)
	var clientID int
	if pl.UserExistsByInput(input) ||
		pl.GetUserIDByInput(input, &clientID) {
		return pl.Err, nil
	}

	tokenPair, session := func(userId int) (tokenPair *models.FreshTokenPair, session *models.RefreshSession) {
		refreshToken := gotp.RandomSecret(rules.RefreshTokenLength)
		session = &models.RefreshSession{
			RefreshToken: refreshToken,
			UserAgent:    "blank agent",
			Exp:          time.Now().Unix() + int64(time.Hour),
			CreatedAt:    time.Now().Unix(),
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			ExpiresAt: session.Exp,
			Subject:   strconv.Itoa(userId),
		}).SignedString([]byte(os.Getenv("SECRET_SIGNING_KEY")))
		if err != nil {
			panic(err)
		}
		tokenPair = &models.FreshTokenPair{
			AccessToken:  token,
			RefreshToken: refreshToken,
		}
		return
	}(clientID)
	sessionsCount, err := r.Services.Repos.Auth.CreateNewUserRefreshSession(clientID, session)
	if err != nil {
		return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
	}
	if sessionsCount > 5 {
		err = r.Services.Repos.Auth.DeleteOldestSession(clientID)
		if err != nil {
			return resp.Error(resp.ErrInternalServerError, "внутренняя ошибка сервера"), nil
		}
	}
	return model.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

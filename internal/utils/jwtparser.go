package utils

import (
	"context"
	"github.com/robbert229/jwt"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/res"
)

type TokenData struct {
	UserID    int   `json:"userid"`
	ExpiresAt int64 `json:"exp"`
}

func ParseToken(tokenString string, secretKey string) (*TokenData, error) {

	var (
		userID     int
		expiresAt  int64
		data       *TokenData
		err        error
		claims     *jwt.Claims
		_userID    interface{}
		_expiresAt interface{}
		fuserID    float64
		fexpiresAt float64
		ok         bool
		algorithm  jwt.Algorithm
	)

	algorithm = jwt.HmacSha256(secretKey)
	if err := algorithm.Validate(tokenString); err != nil {
		goto handleError
	}

	claims, err = algorithm.Decode(tokenString)
	if err != nil {
		goto handleError
	}

	_userID, err = claims.Get("userid")
	if err != nil {
		goto handleError
	}
	_expiresAt, err = claims.Get("exp")
	if err != nil {
		goto handleError
	}

	fuserID, ok = _userID.(float64)
	if !ok {
		err = cerrors.New("token not contain userid")
		goto handleError
	}
	fexpiresAt, ok = _expiresAt.(float64)
	if !ok {
		err = cerrors.New("token not contain exp")
		goto handleError
	}

	userID = int(fuserID)
	expiresAt = int64(fexpiresAt)

	data = &TokenData{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	return data, nil

handleError:
	return nil, cerrors.Wrap(err, "не удалось распарсить токен")
}

func GenerateToken(data *TokenData, secretKey string) (string, error) {
	algorithm := jwt.HmacSha256(secretKey)

	claims := jwt.NewClaim()
	claims.Set("userid", data.UserID)
	claims.Set("exp", data.ExpiresAt)

	token, err := algorithm.Encode(claims)

	if err != nil {
		return "", cerrors.Wrap(err, "не удалось сгенерировать токен")
	}

	return token, nil
}

func GetAuthDataFromCtx(ctx context.Context) (authData *TokenData) {
	data, ok := ctx.Value(res.CtxAuthData).(*TokenData)
	if !ok {
		println("GetAuthDataFromCtx: не удалось найти CtxAuthData в контексте")
	}
	return data
}

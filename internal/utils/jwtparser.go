package utils

import (
	"errors"
	"github.com/robbert229/jwt"
)

type TokenData struct {
	UserID    int   `json:"userid"`
	ExpiresAt int64 `json:"exp"`
}

func ParseToken(tokenString string, secretKey string) (*TokenData, error) {
	println("token:", tokenString) // debug

	algorithm := jwt.HmacSha256(secretKey)
	if err := algorithm.Validate(tokenString); err != nil {
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}

	claims, err := algorithm.Decode(tokenString)
	if err != nil {
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}

	_userID, err := claims.Get("userid")
	if err != nil {
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}
	_expiresAt, err := claims.Get("exp")
	if err != nil {
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}

	fuserID, ok := _userID.(float64)
	if !ok {
		err = errors.New("token not contain userid")
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}
	fexpiresAt, ok := _expiresAt.(float64)
	if !ok {
		err = errors.New("token not contain exp")
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}
	userID := int(fuserID)
	expiresAt := int64(fexpiresAt)
	println(userID)
	println(expiresAt)
	data := &TokenData{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	if err != nil {
		println("ParseToken:", err.Error()) // debug
		return nil, err
	}
	return data, err
}

func GenerateToken(data *TokenData, secretKey string) (string, error) {
	algorithm := jwt.HmacSha256(secretKey)

	claims := jwt.NewClaim()
	claims.Set("userid", data.UserID)
	claims.Set("exp", data.ExpiresAt)

	token, err := algorithm.Encode(claims)

	if err != nil {
		println("GenerateToken:", err.Error()) // debug
		return "", err
	}

	return token, nil
}
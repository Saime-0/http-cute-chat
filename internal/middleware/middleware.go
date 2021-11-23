package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var secretkey = os.Getenv("SECRET_SIGNING_KEY")

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId int
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) == 2 {
			jwtToken := authHeader[1]
			token, _ := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretkey), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userId, _ = strconv.Atoi(claims["sub"].(string))
			}
		}
		ctx := context.WithValue(r.Context(), rules.UserIDFromToken, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

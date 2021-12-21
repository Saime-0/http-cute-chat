package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var secretkey = os.Getenv("SECRET_SIGNING_KEY")

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			expiresAt int64
			userId    int
		)
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
				expiresAt = int64(claims["exp"].(float64))
				if expiresAt >= time.Now().Unix() { // handle expiresAt
					userId, _ = strconv.Atoi(claims["sub"].(string))
				}
			}
		}
		ctx := context.WithValue(r.Context(), rules.UserIDFromToken, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), rules.UserAgentFromHeaders, r.UserAgent())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) Status() int {
	return rw.status
}
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}
func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		log.Println(
			"status", wrapped.status,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"duration", time.Since(start),
		)
	}

	return http.HandlerFunc(fn)
}

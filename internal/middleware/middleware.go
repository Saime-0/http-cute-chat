package middleware

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type chain struct {
	r   *http.Request
	cfg *config.Config
	//ctx context.Context
}

func ChainShip(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &chain{
				r:   r,
				cfg: cfg,
			}

			if r.Header.Get("Sec-Websocket-Protocol") != "graphql-ws" {
				c.checkAuth().getUserAgent()
			} else {
				println("WebsocketExeption working!") // debug
			}

			next.ServeHTTP(w, c.r)
		})
	}
}

// deprecated
func WebsocketExeption() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Sec-Websocket-Protocol") == "graphql-ws" {
				println("WebsocketExeption working!") // debug
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (c *chain) checkAuth() *chain {
	println("CheckAuth start!") // debug

	println(c.r.Header.Get("Authorization")) // debug
	ctx, err := auth(c.r.Context(), c.cfg, c.r.Header.Get("Authorization"))
	if err != nil {
		println("CheckAuth:", err.Error())
	}
	c.r = c.r.WithContext(ctx)
	return c
}

func (c *chain) getUserAgent() *chain {
	println("GetUserAgent start!") // debug

	c.r = c.r.WithContext(context.WithValue(c.r.Context(), rules.UserAgentFromHeaders, c.r.UserAgent()))
	return c
}

func Logging(cfg *config.Config) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, _ := json.MarshalIndent(r.Header, "", " ")
			fmt.Printf("header: %s\n", string(b)) // debug

			println("Logging start!") // debug

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			log.Println(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		})
	}
}

type responseWriter struct {
	http.Hijacker
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
	}

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

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func WebsocketInitFunc(cfg *config.Config) func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {

	return func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
		println("INIT FUNC") // debug
		ctx, err := auth(ctx, cfg, initPayload.Authorization())
		if err != nil {
			println("WebsocketInitFunc:", err.Error())
			return nil, err
		}
		return ctx, nil
	}
}

func auth(ctx context.Context, cfg *config.Config, authHeader string) (context.Context, error) {
	var (
		userId int
		err    error
		data   *utils.TokenData
	)
	token := strings.Split(authHeader, "Bearer ")
	if len(token) == 2 {
		data, err = utils.ParseToken(
			token[1],
			cfg.SecretKey,
		)
		if err == nil && data.ExpiresAt >= time.Now().Unix() {
			userId = data.UserID
		}
		fmt.Printf("%#v, %#v\n", data, err)
	}
	return context.WithValue(ctx, rules.UserIDFromToken, userId), err
}

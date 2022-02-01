package middleware

import (
	"bufio"
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/res"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"net"
	"net/http"
	"strings"
	"time"
)

type chain struct {
	r      *http.Request
	cfg    *config.Config
	logger *clog.Clog
	hlr    *healer.Healer
}

func ChainShip(cfg *config.Config, logger *clog.Clog, hlr *healer.Healer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := &chain{
				r:      r,
				cfg:    cfg,
				logger: logger,
				hlr:    hlr,
			}

			if r.Header.Get("Sec-Websocket-Protocol") != "graphql-ws" {
				c.checkAuth().getUserAgent()
			} else {
				err := c.logger.Debug("connection switched to websocket!")
				c.hlr.MonitorLogger(err)
			}

			next.ServeHTTP(w, c.r)
		})
	}
}

func (c *chain) checkAuth() *chain {
	ctx, err := auth(c.r.Context(), c.cfg, c.r.Header.Get("Authorization"))
	if err != nil {
		c.logger.Alert(err)
	}
	c.r = c.r.WithContext(ctx)
	return c
}

func (c *chain) getUserAgent() *chain {
	//println("(chain)getUserAgent start!") // debug
	c.r = c.r.WithContext(context.WithValue(
		c.r.Context(),
		res.CtxUserAgent,
		c.r.UserAgent(),
	))
	return c
}

func Logging(cfg *config.Config, logger *clog.Clog, hlr *healer.Healer) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			err := logger.Debug(
				bson.M{
					"status":   wrapped.status,
					"method":   r.Method,
					"path":     r.URL.EscapedPath(),
					"duration": time.Since(start).String(),
				},
			)
			hlr.MonitorLogger(err)
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

		ctx, err := auth(ctx, cfg, initPayload.Authorization())
		if err != nil {
			println("WebsocketInitFunc:", err.Error()) // debug
			return nil, err
		}
		return ctx, nil
	}
}

func auth(ctx context.Context, cfg *config.Config, authHeader string) (context.Context, error) {
	var (
		err  error
		data *utils.TokenData
	)
	token := strings.Split(authHeader, "Bearer ")
	if len(token) == 2 {
		data, err = utils.ParseToken(
			token[1],
			cfg.SecretKey,
		)
	}
	return context.WithValue(ctx, res.CtxAuthData, data), err
}

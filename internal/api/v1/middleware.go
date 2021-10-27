package v1

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/saime-0/http-cute-chat/internal/api/responder"
	"github.com/saime-0/http-cute-chat/internal/api/rules"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responder.Writer{
			ResponseWriter: w,
			Code:           http.StatusOK,
		}
		next.ServeHTTP(rw, r)

		log.Printf(
			"completed with %d %s in %v\n",
			rw.Code,
			http.StatusText(rw.Code),
			time.Since(start),
		)
	})
}

// todo, create cfg package
var SECRETKEY = os.Getenv("SECRET_SIGNING_KEY")

func checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(SECRETKEY), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				user_id, _ := strconv.Atoi(claims["sub"].(string))
				ctx := context.WithValue(r.Context(), rules.UserIDFromToken, user_id)
				// Access context values in handlers like this
				// props, _ := r.Context().Value("props").(jwt.MapClaims)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}

// ? Может вынести в отдельный файл
type pipeline struct {
	// c context.Context
	w http.ResponseWriter
	r *http.Request
	h *Handler
}

func initPipeline(w http.ResponseWriter, r *http.Request, h *Handler) *pipeline {
	return &pipeline{
		w: w,
		r: r,
		h: h,
	}
}

func (p *pipeline) finalInspectionDatabase(err error) {
	switch {
	case err == sql.ErrNoRows:
		responder.Error(p.w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		panic(err)

	case err != nil:
		responder.Error(p.w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)
	}
}

func (p *pipeline) parseUserDomainFromRequest() (user_domain string) {
	user_domain = mux.Vars(p.r)["user-domain"]
	if !validateDomain(user_domain) {
		err := rules.ErrInvalidValue
		responder.Error(p.w, http.StatusBadRequest, err)

		panic(err)
	}

	return
}

func (p *pipeline) parseUserIDFromRequest() (user_id int) {
	user_id, err := strconv.Atoi(mux.Vars(p.r)["user-id"])
	if err != nil {
		err := rules.ErrInvalidValue
		responder.Error(p.w, http.StatusBadRequest, err)

		panic(err)
	}

	return
}

func (p *pipeline) inspectUserExistsByDomain(user_domain string) {
	if !p.h.Services.Repos.Users.UserExistsByDomain(user_domain) {
		err := rules.ErrUserNotFound
		responder.Error(p.w, http.StatusBadRequest, err)

		log.Panic(err)
	}

	return
}

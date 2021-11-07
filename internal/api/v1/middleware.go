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

func finalInspectionDatabase(w http.ResponseWriter, err error) {
	switch {
	case err == sql.ErrNoRows:
		responder.Error(w, http.StatusInternalServerError, rules.ErrDataRetrieved)
		panic(err)

	case err != nil:
		responder.Error(w, http.StatusInternalServerError, rules.ErrAccessingDatabase)
		panic(err)
	}
}

func parseOffsetFromQuery(w http.ResponseWriter, r *http.Request) (offset int, ok bool) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil && r.URL.Query().Get("offset") != "" {
		responder.Error(w, http.StatusBadRequest, rules.ErrInvalidValue)

		return
	}

	if offset < 0 {
		responder.Error(w, http.StatusBadRequest, rules.ErrOutOfRange)

		return
	}
	return offset, true
}

// мб скомбинировать pipline и обычные проверки?

// todo
func UserHaveAccessToManageRoom(user_id int, room_id int) (have bool) {
	// get user role
	return // role.RoomManage == true && (room_id == nil || room_id == role.RoomID)
}
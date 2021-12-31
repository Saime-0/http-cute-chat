package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/directive"
	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/resolver"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/middleware"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var configpath string

func init() {
	flag.StringVar(&configpath, "cfg", "cute-config.toml", "path to configure config")
}

func main() {
	flag.Parse()
	cfg := config.NewConfig(configpath)

	// init database
	db, err := initDB(cfg)
	if err != nil {
		println("Initialization database failure:", err.Error())
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Print("error when closing the database")
		}
	}(db)

	services := service.NewServices(db)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			Services: services,
			Config:   cfg,
			Piper:    piper.NewPipeline(services.Repos),
		},
		Directives: generated.DirectiveRoot{
			IsAuth:        directive.IsAuth,
			InputUnion:    directive.InputUnion,
			InputLeastOne: directive.InputLeastOne,
		},
	}))

	router := mux.NewRouter()
	mw := middleware.Setup(cfg)
	router.Use(
		mw.Logging,
		mw.CheckAuth,
		mw.GetUserAgent,
	)

	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 1,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  0, // reused buffers
			WriteBufferSize: 0,
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			println("INIT FUNC") // debug
			var (
				expiresAt int64
				userId    int
			)
			authHeader := strings.Split(initPayload.Authorization(), "Bearer ")
			if len(authHeader) == 2 {
				jwtToken := authHeader[1]
				token, _ := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}

					return []byte(cfg.SecretKey), nil
				})

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					expiresAt = int64(claims["exp"].(float64))
					if expiresAt >= time.Now().Unix() { // handle expiresAt
						userId, _ = strconv.Atoi(claims["sub"].(string))
					}
				}
			}
			ctx = context.WithValue(ctx, rules.UserIDFromToken, userId)
			return ctx, nil
		},
	})
	srv.Use(extension.Introspection{})

	//c := cors.Default().Handler(router)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, router))
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DbName,
	)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil

}

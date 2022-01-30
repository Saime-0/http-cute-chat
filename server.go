package main

import (
	"flag"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/directive"
	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/resolver"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/middleware"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/service"
	"github.com/saime-0/http-cute-chat/internal/store"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"net/http"
	"time"
)

var configpath string

func init() {
	flag.StringVar(&configpath, "cfg", "cute-config.toml", "path to configure config")
}

func main() {
	var err error
	flag.Parse()
	cfg := config.NewConfig(configpath)

	// init logger
	logger, err := clog.NewClog(cfg)
	if err != nil {
		panic(err)
	}

	// init database
	db, err := store.InitDB(cfg)
	if err != nil {
		logger.Emergency(err.Error())
		return
	}
	defer db.Close()

	// init services
	services := service.NewServices(db, cfg, logger)

	// server handler
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
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
		Complexity: *utils.MatchComplexity(),
	}))

	// init router and middlewares
	router := mux.NewRouter()
	router.Use(
		middleware.Logging(cfg, logger),
		middleware.ChainShip(cfg),
	)

	// configure available request methods
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Minute,
			CheckOrigin: func(r *http.Request) bool {
				// todo we are already checking for CORS
				return true
			},
			EnableCompression: true,
			ReadBufferSize:    0, // reused buffers
			WriteBufferSize:   0,
		},
		InitFunc: middleware.WebsocketInitFunc(cfg),
	})

	// server capabilities
	srv.Use(extension.Introspection{})
	srv.Use(extension.FixedComplexityLimit(cfg.QueryComplexityLimit))

	//c := cors.Default().Handler(router)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	err = logger.Info(fmt.Sprintf("Server started on %s port", cfg.AppPort))
	if err != nil {
		panic(err)
	}
	logger.Alert(http.ListenAndServe(":"+cfg.AppPort, router).Error())
}

package main

import (
	"flag"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/directive"
	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/resolver"
	"github.com/saime-0/http-cute-chat/internal/cache"
	"github.com/saime-0/http-cute-chat/internal/clog"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/email"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/middleware"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/rules"
	"github.com/saime-0/http-cute-chat/internal/service"
	"github.com/saime-0/http-cute-chat/internal/store"
	"github.com/saime-0/http-cute-chat/internal/subix"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/graphiql"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
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
	logger, err := clog.NewClog(cfg, clog.Multiple)
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

	// init smtp
	newSMTPSender, err := email.NewSMTPSender(
		cfg.SMTP.Author,
		cfg.SMTP.From,
		cfg.SMTP.Passwd,
		cfg.SMTP.Host,
		cfg.SMTP.Port,
	)
	if err != nil {
		logger.Emergency(err.Error())
		return
	}
	// init services
	services := &service.Services{
		Repos:     repository.NewRepositories(db),
		Scheduler: scheduler.NewScheduler(),
		SMTP:      newSMTPSender,
		Cache:     cache.NewCache(),
		Logger:    logger,
	}

	// init healer
	hlr, err := healer.NewHealer(cfg, services.Scheduler, services.Cache, services.Logger)
	if err != nil {
		logger.Emergency(err.Error())
		return
	}

	// init subix
	sbx := subix.NewSubix(services.Repos, services.Scheduler)

	// init resolver
	myResolver := &resolver.Resolver{
		Services: services,
		Config:   cfg,
		Piper:    piper.NewPipeline(services.Repos, hlr),
		Healer:   hlr,
		Subix:    sbx,
	}
	err = myResolver.RegularSchedule(rules.DurationOfScheduleInterval)
	if err != nil {
		logger.Emergency(err.Error())
		return
	}

	// server handler
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: myResolver,
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
		middleware.Logging(cfg, logger, hlr),
		middleware.ChainShip(cfg, logger, hlr),
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

	// handlers
	router.Handle("/", graphiql.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	err = logger.Info(fmt.Sprintf("Server started on %s port", cfg.AppPort))
	if err != nil {
		logger.Emergency(err.Error())
		return
	}
	logger.Alert(http.ListenAndServe(":"+cfg.AppPort, router).Error())
}

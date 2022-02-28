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
	"github.com/saime-0/http-cute-chat/internal/cdl"
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/email"
	"github.com/saime-0/http-cute-chat/internal/healer"
	"github.com/saime-0/http-cute-chat/internal/middleware"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/repository"
	"github.com/saime-0/http-cute-chat/internal/service"
	"github.com/saime-0/http-cute-chat/internal/store"
	"github.com/saime-0/http-cute-chat/internal/subix"
	"github.com/saime-0/http-cute-chat/internal/utils"
	"github.com/saime-0/http-cute-chat/pkg/graphiql"
	"github.com/saime-0/http-cute-chat/pkg/scheduler"
	"net/http"
	"os"
	"time"
)

var configpath string

func init() {
	flag.StringVar(&configpath, "cfg", "cute-config.toml", "path to configure config")
}

func main() {
	var err error
	flag.Parse()
	cfg, err := config.NewConfig2(configpath)
	if err != nil {
		panic(err)
	}

	newSched := scheduler.NewScheduler()
	newCache := cache.NewCache()

	// init healer
	hlr, err := healer.NewHealer(cfg, newSched, newCache)
	if err != nil {
		panic(cerrors.Wrap(err, "ошибка инициализации компонента восстановления"))
	}
	//hlr.Emergency("test emegrency")

	// init database
	db, err := store.InitDB(cfg)
	if err != nil {
		hlr.Emergency(cerrors.Wrap(err, "ошибка соединения с бд").Error())
		os.Exit(69)
	}
	defer db.Close()

	// init smtp
	newSMTPSender, err := email.NewSMTPSender(
		*cfg.SMTPing.SMTPAuthor,
		cfg.SmtpEmailLogin,
		cfg.SmtpEmailPasswd,
		cfg.SmtpHost,
		*cfg.SMTPing.SMTPPort,
	)
	if err != nil {
		hlr.Emergency(cerrors.Wrap(err, "ошибка создания компонента SMTP").Error())
		os.Exit(69)
	}
	// init services
	services := &service.Services{
		Repos:     repository.NewRepositories(db),
		Scheduler: newSched,
		SMTP:      newSMTPSender,
		Cache:     newCache,
	}

	// init subix
	sbx := subix.NewSubix(services.Repos, services.Scheduler)

	// init dataloader
	dataloader := cdl.NewDataloader(time.Millisecond*5, 100, db, hlr)

	// init resolver
	myResolver := &resolver.Resolver{
		Services:   services,
		Config:     cfg,
		Piper:      piper.NewPipeline(cfg, services.Repos, hlr, dataloader),
		Healer:     hlr,
		Subix:      sbx,
		Dataloader: dataloader,
	}
	err = myResolver.RegularSchedule(*cfg.DurationOfScheduleInterval)
	if err != nil {
		hlr.Emergency(cerrors.Wrap(err, "ошибка компонента очистки неактуальных записей базы данных").Error())
		os.Exit(69)
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
		middleware.InitNode(myResolver.Piper, hlr),
		middleware.ChainShip(cfg, hlr),
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
	srv.Use(extension.FixedComplexityLimit(*cfg.QueryComplexityLimit))

	// handlers
	router.Handle("/", graphiql.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	hlr.Info(fmt.Sprintf("Server started on %s port", *cfg.ApplicationPort))
	err = http.ListenAndServe(":"+*cfg.ApplicationPort, router)
	if err != nil {
		hlr.Emergency(cerrors.Wrap(err, "ListenAndServe failure").Error())
		os.Exit(69)
	}
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/directive"
	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/resolver"
	"github.com/saime-0/http-cute-chat/internal/config"
	"github.com/saime-0/http-cute-chat/internal/middleware"
	"github.com/saime-0/http-cute-chat/internal/piper"
	"github.com/saime-0/http-cute-chat/internal/service"
	"log"
	"net/http"
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
	// srv.AddTransport(transport.POST{})
	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 0,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},

			EnableCompression: true,
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
		},
	})

	router := mux.NewRouter()
	mw := middleware.Setup(cfg)
	router.Use(
		mw.Logging,
		mw.CheckAuth,
		mw.GetUserAgent,
	)

	//c := cors.New(cors.Options{
	//	AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
	//	AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return true
	//	},
	//})
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

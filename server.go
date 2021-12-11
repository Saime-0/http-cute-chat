package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/graph/directive"
	"github.com/saime-0/http-cute-chat/graph/model"
	"github.com/saime-0/http-cute-chat/internal/middleware"
	"github.com/saime-0/http-cute-chat/internal/service"
	"log"
	"net/http"
	"os"
	//
	"github.com/saime-0/http-cute-chat/graph/generated"
	"github.com/saime-0/http-cute-chat/graph/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

const (
	host     = "localhost"
	dbPort   = 5432
	user     = "postgres"
	password = "7050"
	dbname   = "chat_db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	// init database
	db, err := initDB()
	if err != nil {
		panic(err)
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
		},
		Directives: generated.DirectiveRoot{
			HasChar: func(ctx context.Context, obj interface{}, next graphql.Resolver, char []*model.CharType) (res interface{}, err error) {
				return next(ctx)
			},
			IsAuth:     directive.IsAuth,
			InputUnion: directive.InputUnion,
		},
	}))
	router := mux.NewRouter()
	router.Use(middleware.CheckAuth)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initDB() (*sql.DB, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbname)

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

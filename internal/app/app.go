package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/saime-0/http-cute-chat/internal/api"
	"github.com/saime-0/http-cute-chat/internal/service"
)

const (
	serverAddr = ":8081"
	host       = "localhost"
	dbPort     = 5432
	user       = "postgres"
	password   = "7050"
	dbname     = "chat_db"
)

type ApiServer struct {
	httpServer *http.Server
	store      *Store
}

func newApiServer(db *sql.DB, handler http.Handler) *ApiServer {
	a := &ApiServer{
		httpServer: NewHttpServer(serverAddr, handler),
		store:      NewStore(db),
	}
	return a
}

func Run() {
	log.Println("Run")
	start := time.Now()

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

	// new services
	services := service.NewServices(db)

	handler := api.NewGeneralHandler(services)

	// create api server struct
	a := newApiServer(db, handler.Init())

	// run HTTP Server
	go func() {
		log.Println("server startup was successful in", time.Since(start))
		if err = a.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	start = time.Now()
	a.Stop(context.TODO())
	log.Println(log.Prefix(), "Graceful Shutdown in", time.Since(start))

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

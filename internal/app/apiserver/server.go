package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paramonies/sber-rest-api/internal/app/controller"
	"github.com/paramonies/sber-rest-api/internal/app/repository"
	"github.com/paramonies/sber-rest-api/internal/app/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Server provides an http.Server
type Server struct {
	*http.Server
	DB *sqlx.DB
}

func NewServer(config *Config) (*Server, error) {
	db, err := newDB(config)
	if err != nil {
		return &Server{}, err
	}

	api := newAPI(db)

	srv := http.Server{
		Addr:         config.SrvHost + ":" + config.SrvPort,
		Handler:      api,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{&srv, db}, nil
}

// Start runs ListenAndServe on the http.Server with graceful shutdown
func (srv *Server) Start() {
	log.Println("Starting server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Could not listen on addr", srv.Addr, "reason ", err.Error())
		}
	}()
	log.Println("Server is ready to handle requests addr", srv.Addr)
	srv.gracefulShutdown()
}

func (srv *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Server is shutting down, reason", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Could not gracefully shutdown the server", err.Error())
	}
	log.Println("Server stopped")

	if err := srv.DB.Close(); err != nil {
		log.Fatal("Could not close DataBase", err.Error())
	}
	log.Println("DataBase connection closed")
}

func newDB(config *Config) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newAPI(db *sqlx.DB) *gin.Engine {
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	controller := controller.NewController(service)

	return controller.InitRoutes()
}

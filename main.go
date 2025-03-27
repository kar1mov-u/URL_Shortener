package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shortener/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	//Creating DB connectcion and SQLC
	DB_URL := os.Getenv("DB_URL")
	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		panic(err)
	}
	sqlc := database.New(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Configuring/Starting serves
	worker := Worker{DB: sqlc, CTX: ctx}

	//TTL Deleter
	go func() {
		for {
			worker.ttlDeleter()
			time.Sleep(time.Hour * 5)
		}
	}()

	//Configuring server
	cfg := Config{DB: sqlc}
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	//Web-endponts
	router.Get("/", cfg.home)
	router.Get("/redirect/{hashUrl}", cfg.redirect)

	//API end-points
	router.Post("/urls/shorten", cfg.createUrl)
	router.Get("/urls/convert", cfg.convert)

	//handling shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		log.Println("Killing workers")
		cancel()
		os.Exit(0)
	}()

	log.Println("Starting server on port: 8080")
	http.ListenAndServe(":8080", router)
}

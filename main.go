package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"shortener/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	DB_URL := os.Getenv("DB_URL")
	//craete connection to db
	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		panic(err)
	}
	//connection to sqlc
	sqlc := database.New(conn)

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

	log.Println("Starting server on port: 8080")
	http.ListenAndServe(":8080", router)
}

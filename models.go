package main

import "shortener/internal/database"

type Config struct {
	DB *database.Queries
}
type URL struct {
	Original string `json:"original"`
}

type ErrorResp struct {
	Error string `json:"err"`
}

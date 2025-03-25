package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"shortener/internal/database"
)

func (cfg *Config) createUrl(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Original string `json:"original"`
	}
	data := Request{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithErr(w, 500, err.Error())
	}

	//check if it exists, if yes return it
	dbURL, err := cfg.DB.GetUrlbyOrig(r.Context(), data.Original)
	if err != nil {
		//checking if this is not ErrNoRows, if its we can continue and generate a hash
		if !errors.Is(err, sql.ErrNoRows) {
			respondWithErr(w, 500, err.Error())
		}
	}
	//if its not empty means we have already hashed url in DB
	if dbURL != "" {
		respondWithJson(w, 200, map[string]string{"short_url": dbURL})
		return
	}

	//if not create hash
	hashURL := generateHash(data.Original)
	err = cfg.DB.CreateEntry(r.Context(), database.CreateEntryParams{
		OriginalUrl: data.Original,
		HashedUrl:   hashURL,
	})
	if err != nil {
		respondWithErr(w, 500, err.Error())
		return
	}

	respondWithJson(w, 200, map[string]string{"short_url": dbURL})
}

func (cfg *Config) convert(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Hash string `json:"hash"`
	}
	data := Request{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		respondWithErr(w, 500, err.Error())
	}

	url, err := cfg.DB.GetUrlbyHash(r.Context(), data.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithErr(w, 404, "No such entry")
			return
		} else {
			respondWithErr(w, 500, err.Error())
			return
		}
	}
	respondWithJson(w, 200, map[string]string{"url": url})
}

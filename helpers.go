package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func respondWithErr(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func generateHash(origUrl string) string {
	hash := sha256.Sum256([]byte(origUrl))
	hashString := hex.EncodeToString(hash[:])
	return hashString[:7]
}

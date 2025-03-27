package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var (
	homeTmpl *template.Template
)

func init() {
	var err error
	homeTmpl, err = template.ParseFiles("./html/index.html")
	if err != nil {
		panic(err)
	}
}

func (cfg *Config) home(w http.ResponseWriter, r *http.Request) {
	err := homeTmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (cfg *Config) redirect(w http.ResponseWriter, r *http.Request) {

	hashUrl := chi.URLParam(r, "hashUrl")
	if len(hashUrl) != 7 {
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}

	//make api call
	payload := fmt.Sprintf(`{"hash":"%v"}`, hashUrl)
	req, err := http.NewRequest("GET", "http://localhost:8080/urls/convert", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on making API call on handler redirect")
		http.Error(w, "Something went wrong", 500)
		return
	}
	defer resp.Body.Close()

	//Check for errors
	if resp.StatusCode != 200 {
		errResp := ErrorResp{}
		http.Error(w, errResp.Error, resp.StatusCode)
		return
	}

	//extract url
	type Response struct {
		Url string `json:"url"`
	}
	respBody := Response{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Println("Failed to parse json on handler redirect")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, respBody.Url, http.StatusPermanentRedirect)

}

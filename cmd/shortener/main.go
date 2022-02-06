package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

var storage map[string]string

var urlID int

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		link := string(b)
		shortURL := strconv.Itoa(urlID)
		urlID++
		storage[shortURL] = link
		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + string(shortURL)))
	case http.MethodGet:
		path := r.URL.Path[1:]
		link := storage[path]
		w.Header().Add("Location", link)
		w.WriteHeader(307)
	default:
		w.WriteHeader(400)
		w.Write([]byte("Something wrong"))
	}
}

func main() {
	storage = make(map[string]string)
	urlID = 1
	http.HandleFunc("/", ShortURLHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

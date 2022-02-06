package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type OriginalUrl struct {
	Location string `json:"Location"`
}

var storage map[string]string

var urlId int

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		link := string(b)
		shortUrl := strconv.Itoa(urlId)
		urlId++
		storage[shortUrl] = link
		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + string(shortUrl)))
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
	urlId = 1
	http.HandleFunc("/", ShortUrlHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

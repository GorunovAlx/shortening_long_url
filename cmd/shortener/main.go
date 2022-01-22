package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type OriginalUrl struct {
	Location string `json:"Location"`
}

var storage map[string]string

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		link := string(b)
		shortUrl := GenerateShortLink(link)
		storage[shortUrl] = link
		w.WriteHeader(201)
		w.Write([]byte("http://localhost:8080/" + shortUrl))
	case http.MethodGet:
		path := r.URL.Path[1:]
		link := storage[path]
		w.Header().Set("content-type", "application/json")
		originalUrl := OriginalUrl{link}
		resp, err := json.Marshal(originalUrl)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(resp)
		w.WriteHeader(307)
	default:
		w.WriteHeader(400)
	}
}

func main() {
	storage = make(map[string]string)
	// маршрутизация запросов обработчику
	http.HandleFunc("/", ShortUrlHandler)
	// запуск сервера с адресом localhost, порт 8080
	http.ListenAndServe(":8080", nil)
}

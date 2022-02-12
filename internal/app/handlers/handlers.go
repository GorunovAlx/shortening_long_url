package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func CreateShortURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Incorrect request"))
			return
		}
		if string(b) == "" {
			w.WriteHeader(400)
			w.Write([]byte("Incorrect request"))
			return
		}

		shortURL, err := urlStorage.CreateShortURL(string(b))
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(201)
		urlNumber := strconv.Itoa(shortURL)
		w.Write([]byte("http://localhost:8080/" + urlNumber))
	}
}

func GetInitialLinkHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		if shortURL == "" {
			w.WriteHeader(400)
			w.Write([]byte("short url was not sent"))
			return
		}

		url, err := strconv.Atoi(shortURL)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		link, err := urlStorage.GetInitialLink(url)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Location", link)
		w.WriteHeader(307)
	}
}

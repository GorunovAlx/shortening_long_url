package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func ShortenerURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
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
		case http.MethodGet:
			path, err := strconv.Atoi(r.URL.Path[1:])
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				return
			}

			link, err := urlStorage.GetInitialLink(path)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
				return
			}

			w.Header().Add("Location", link)
			w.WriteHeader(307)
		default:
			w.WriteHeader(400)
			w.Write([]byte("Something wrong"))
		}
	}
}

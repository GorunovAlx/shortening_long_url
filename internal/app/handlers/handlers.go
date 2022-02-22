package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

// Handler is a structure containing the type of chi.Mux
// and ShortURLRepo interface from storage package.
type Handler struct {
	*chi.Mux
	Repo storage.ShortURLRepo
}

// NewHandler returns a newly initialized Handler object that implements
// the ShortURLRepo interface.
func NewHandler(repo storage.ShortURLRepo) *Handler {
	h := &Handler{
		Mux:  chi.NewMux(),
		Repo: repo,
	}

	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	h.Post("/", CreateShortURLHandler(repo))
	h.Get("/{shortURL}", GetInitialLinkHandler(repo))
	h.Post("/api/shorten", CreateShortURLJSONHandler(repo))

	return h
}

func RegisterRoutes(repo storage.ShortURLRepo) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", CreateShortURLHandler(repo))
		r.Get("/{shortURL}", GetInitialLinkHandler(repo))
		r.Post("/api/shorten", CreateShortURLJSONHandler(repo))
	})

	return r
}

func CreateShortURLJSONHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url storage.ShortURL
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isURL := valid.IsURL(url.InitialLink)
		if !isURL {
			w.WriteHeader(400)
			w.Write([]byte("Incorrect link"))
			return
		}

		shortURL, err := urlStorage.CreateShortURL(url.InitialLink)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		res := storage.ShortURL{
			ShortLink: shortURL,
		}
		resp, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}

// CreateShortURLHandler returns a http.HandlerFunc that processes the body of the request
// which contains URL and returns a shortened link.
func CreateShortURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Incorrect request"))
			return
		}
		if len(b) == 0 {
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

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + shortURL))
	}
}

// GetInitialLinkHandler returns a http.HandlerFunc that takes shortURL parameter
// containing a short url and returns the initial link in the location header.
func GetInitialLinkHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		if shortURL == "" {
			w.WriteHeader(400)
			w.Write([]byte("short url was not sent"))
			return
		}

		link, err := urlStorage.GetInitialLink(shortURL)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

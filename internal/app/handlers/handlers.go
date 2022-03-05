package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/caarlos0/env/v6"
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

//
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}

var Cfg Config

func Init() {
	err := env.Parse(&Cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(Cfg)
}

// NewHandler returns a newly initialized Handler object that implements
// the ShortURLRepo interface.
func NewHandler(repo storage.ShortURLRepo) *Handler {
	Init()
	h := &Handler{
		Mux:  chi.NewMux(),
		Repo: repo,
	}

	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	h.Use(middleware.Timeout(60 * time.Second))

	h.Get("/{shortURL}", GetInitialLinkHandler(repo))
	h.Post("/", CreateShortURLHandler(repo))
	h.Post("/api/shorten", CreateShortURLJSONHandler(repo))

	return h
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
		shortURL = Cfg.BaseURL + shortURL
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
		w.Write([]byte(Cfg.BaseURL + shortURL))
	}
}

// GetInitialLinkHandler returns a http.HandlerFunc that takes shortURL parameter
// containing a short url and returns the initial link in the location header.
func GetInitialLinkHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		log.Println("Get link on path: " + shortURL)
		if shortURL == "" {
			w.WriteHeader(400)
			w.Write([]byte("short url was not sent"))
			return
		}

		link, err := urlStorage.GetInitialLink(shortURL)
		log.Println("Get initial link: " + link)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Location", link)
		w.WriteHeader(307)
	}
}

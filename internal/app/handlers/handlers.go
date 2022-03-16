package handlers

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

// Returns a pointer to a chi.Mux with endpoints:
// Get /{shortURL} returns the initial link from storage by shortened link.
// Post / sends initial link in the body and get shortened link in the response body.
// Post /api/shorten sends json with initial link in the body
// and get json with shortened link in the response body.
func NewRouter(repo storage.ShortURLRepo) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(mGzipWriterHandle)
	r.Use(mGzipReaderHandle)

	r.Get("/{shortURL}", GetInitialLinkHandler(repo))
	r.Post("/", CreateShortURLHandler(repo))
	r.Post("/api/shorten", CreateShortURLJSONHandler(repo))

	return r
}

//
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

//
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

//
func mGzipWriterHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

//
func mGzipReaderHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		var reader io.Reader

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}

		body := io.NopCloser(reader)
		r.Body = body

		w.Header().Set("Accept-Encoding", "gzip")
		next.ServeHTTP(w, r)
	})
}

// Post a json with an initial link in the request and returns a json
// with a shortened link in the response.
func CreateShortURLJSONHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url storage.ShortURL
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		isURL := valid.IsURL(url.InitialLink)
		if !isURL {
			http.Error(w, "Incorrect link", http.StatusBadRequest)
			return
		}

		shortURL, err := urlStorage.CreateShortURL(url.InitialLink)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := storage.ShortURL{
			ShortLink: configs.Cfg.BaseURL + "/" + shortURL,
		}
		resp, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}

// Post initial link in the request and returns a shortened link in the response.
func CreateShortURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(b) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		shortURL, err := urlStorage.CreateShortURL(string(b))
		shortURL = configs.Cfg.BaseURL + "/" + shortURL
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}
}

// GetInitialLinkHandler returns a http.HandlerFunc that takes shortURL parameter
// containing a short url and returns the initial link in the location header.
func GetInitialLinkHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")
		if shortURL == "" {
			http.Error(w, "short url was not sent", http.StatusBadRequest)
			return
		}

		link, err := urlStorage.GetInitialLink(shortURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Location", link)
		w.WriteHeader(307)
	}
}

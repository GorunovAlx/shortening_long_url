package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"io"
	"net/http"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	gen "github.com/GorunovAlx/shortening_long_url/internal/app/generators"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
	"github.com/GorunovAlx/shortening_long_url/internal/app/utils"
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

	r.Use(MiddlewareGzipWriterHandle)
	r.Use(MiddlewareGzipReaderHandle)
	r.Use(MiddlewareAuthUserHandle)

	r.Get("/{shortURL}", GetInitialLinkHandler(repo))
	r.Get("/api/user/urls", GetAllShortURLUserHandler(repo))
	r.Get("/ping", GetPingToDBHandle(repo))
	r.Post("/", CreateShortURLHandler(repo))
	r.Post("/api/shorten", CreateShortURLJSONHandler(repo))
	r.Post("/api/shorten/batch", CreateListShortURLHandler(repo))
	r.Delete("/api/user/urls", DeleteListURLHandler(repo))

	return r
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

		token := r.Context().Value(contextKeyRequestID).(string)
		id, err := gen.GetUserID(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url.UserID = id

		shortURL, err := urlStorage.CreateShortURL(&url)
		shortURL = configs.Cfg.BaseURL + "/" + shortURL
		if err != nil && err != utils.ErrUniqueLink {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "application/json")
		if errors.Is(err, utils.ErrUniqueLink) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		res := storage.ShortURL{
			ShortLink: shortURL,
		}
		resp, e := json.Marshal(res)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		w.Write(resp)
	}
}

// Post initial link in the request and returns a shortened link in the response.
func CreateShortURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		isURL := valid.IsURL(string(b))
		if !isURL {
			http.Error(w, "Incorrect link", http.StatusBadRequest)
			return
		}

		token := r.Context().Value(contextKeyRequestID).(string)
		id, err := gen.GetUserID(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		shortURL := storage.ShortURL{
			InitialLink: string(b),
			UserID:      id,
		}

		shortened, err := urlStorage.CreateShortURL(&shortURL)
		shortened = configs.Cfg.BaseURL + "/" + shortened
		if err != nil && err != utils.ErrUniqueLink {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, utils.ErrUniqueLink) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		w.Write([]byte(shortened))
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

func GetAllShortURLUserHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getCookieByName("user_id", r)
		if token == "" {
			token = r.Context().Value(contextKeyRequestID).(string)
		}
		id, err := gen.GetUserID(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := urlStorage.GetAllShortURLUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp, err := json.MarshalIndent(res, "", " ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func GetPingToDBHandle(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := urlStorage.PingDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func CreateListShortURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var links []storage.ShortURLByUser
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &links)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := urlStorage.CreateListShortURL(links)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
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

type Job struct {
	shortURL string
}

func DeleteListURLHandler(urlStorage storage.ShortURLRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var links []string
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &links)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := getCookieByName("user_id", r)
		if token == "" {
			token = r.Context().Value(contextKeyRequestID).(string)
		}
		id, err := gen.GetUserID(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		workersCount := 3

		jobCh := make(chan *Job)
		for i := 0; i < workersCount; i++ {
			go func() {
				for job := range jobCh {
					urlStorage.DeleteShortURLUser(job.shortURL, id)
				}
			}()
		}

		for i := 0; i < len(links); i++ {
			job := &Job{shortURL: links[i]}
			jobCh <- job
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

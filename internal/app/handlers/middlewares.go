package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	gen "github.com/GorunovAlx/shortening_long_url/internal/app/generators"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type contextKey int

const (
	contextKeyRequestID contextKey = iota
)

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Give a compressed response to a client that supports compressed response processing.
func MiddlewareGzipWriterHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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

// Accept requests in compressed format.
func MiddlewareGzipReaderHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
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

// MiddlewareAuthUserHandle checks if the user's id cookie came in and if so,
// checks for authentication. If the cookie is empty,
// it creates a new user id cookie, sets it and passes it on.
func MiddlewareAuthUserHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDToken := getCookieByName("user_id", r)

		if len(userIDToken) != 0 {
			isAuthentic, err := gen.AuthUserIDToken(userIDToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if isAuthentic {
				ctx := r.Context()
				ctx = context.WithValue(ctx, contextKeyRequestID, userIDToken)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		userIDToken, err := gen.GenerateUserIDToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{
			Name:    "user_id",
			Value:   userIDToken,
			Path:    "/",
			Expires: expiration,
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyRequestID, userIDToken)
		http.SetCookie(w, &cookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getCookieByName(cName string, r *http.Request) string {
	receivedCookie := r.Cookies()
	var value string
	if len(receivedCookie) != 0 {
		for _, cookie := range receivedCookie {
			if cookie.Name == cName {
				value = cookie.Value
			}
		}
	}

	return value
}

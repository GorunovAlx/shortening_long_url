package handlers

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareGzipWriterHandle(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(w.Header().Get("Content-Encoding"), "gzip") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("content not encoded"))
		}
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Add("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	h := MiddlewareGzipWriterHandle(nextHandler)
	h.ServeHTTP(w, request)
	result := w.Result()

	assert.Equal(t, "gzip", result.Header.Get("Content-Encoding"))
}

func TestMiddlewareGzipReaderHandle(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(w.Header().Get("Accept-Encoding"), "gzip") {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	body := []byte("A long time ago in a galaxy far, far away...")
	comData, err := compress(body)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(comData))
	request.Header.Add("Content-Encoding", "gzip")
	w := httptest.NewRecorder()
	h := MiddlewareGzipReaderHandle(nextHandler)
	h.ServeHTTP(w, request)
	result := w.Result()

	assert.Equal(t, "gzip", result.Header.Get("Accept-Encoding"))
}

func TestMiddlewareAuthUserHandle(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getCookieByName("user_id", r)
		if token == "" {
			token = r.Context().Value(contextKeyRequestID).(string)
		}
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("token is not exists"))
		}
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h := MiddlewareAuthUserHandle(nextHandler)
	h.ServeHTTP(w, request)
	result := w.Result()

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEqual(t, "", result.Header.Get("Set-Cookie"))

}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Name = "a-new-hope.txt"
	w.Comment = "an epic space opera by George Lucas"

	_, err := w.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	return b.Bytes(), nil
}

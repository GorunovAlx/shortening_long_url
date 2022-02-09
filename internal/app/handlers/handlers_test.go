package handlers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	id      int
	storage map[int]string
}

func (ms *mockStorage) idInkrement() {
	ms.id += 1
}

func (ms *mockStorage) GetInitialLink(shortLink int) (string, error) {
	link := ms.storage[shortLink]
	if link == "" {
		return "", errors.New("the url with this value does not exist")
	}
	return link, nil
}

func (ms *mockStorage) CreateShortURL(initialLink string) (int, error) {
	ms.storage[ms.id] = initialLink
	defer ms.idInkrement()

	return ms.id, nil
}

func TestShortenerURLHandler(t *testing.T) {
	type want struct {
		statusCode int
		link       string
	}
	tests := []struct {
		name    string
		method  string
		st      map[int]string
		nextId  int
		request string
		body    []byte
		want    want
	}{
		{
			name:    "simple test #1(Post)",
			request: "http://localhost:8080/",
			method:  http.MethodPost,
			st:      make(map[int]string),
			nextId:  1,
			body:    []byte("google.com"),
			want: want{
				statusCode: 201,
				link:       "http://localhost:8080/1",
			},
		},
		{
			name:    "simple test #2(Get)",
			request: "http://localhost:8080/2",
			method:  http.MethodGet,
			st: map[int]string{
				1: "yandex.ru",
				2: "google.com",
				3: "tutu.ru",
			},
			nextId: 4,
			body:   nil,
			want: want{
				statusCode: 307,
				link:       "google.com",
			},
		},
		{
			name:    "simple test #3(Post)",
			request: "http://localhost:8080/",
			method:  http.MethodPost,
			st: map[int]string{
				1: "yandex.ru",
				2: "google.com",
				3: "tutu.ru",
			},
			nextId: 4,
			body:   []byte("yandex.ru"),
			want: want{
				statusCode: 201,
				link:       "http://localhost:8080/4",
			},
		},
		{
			name:    "simple test #4(Get)",
			request: "http://localhost:8080/5",
			method:  http.MethodGet,
			st: map[int]string{
				1: "yandex.ru",
				2: "google.com",
				3: "tutu.ru",
			},
			nextId: 4,
			body:   nil,
			want: want{
				statusCode: 400,
				link:       "",
			},
		},
		{
			name:    "simple test #5(Get)",
			request: "http://localhost:8080/",
			method:  http.MethodGet,
			st: map[int]string{
				1: "yandex.ru",
				2: "google.com",
				3: "tutu.ru",
			},
			nextId: 4,
			body:   nil,
			want: want{
				statusCode: 400,
				link:       "",
			},
		},
		{
			name:    "simple test #6(Post)",
			request: "http://localhost:8080/",
			method:  http.MethodPost,
			st: map[int]string{
				1: "yandex.ru",
				2: "google.com",
				3: "tutu.ru",
			},
			nextId: 4,
			body:   nil,
			want: want{
				statusCode: 400,
				link:       "Incorrect request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mockStorage{
				id:      tt.nextId,
				storage: tt.st,
			}
			request := httptest.NewRequest(tt.method, tt.request, bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenerURLHandler(&ms))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			var link []byte
			var location string
			var err error
			if tt.method == http.MethodPost {
				link, err = ioutil.ReadAll(result.Body)
			} else {
				location = result.Header.Get("Location")
			}

			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			if tt.method == http.MethodPost {
				assert.Equal(t, tt.want.link, string(link))
			} else {
				assert.Equal(t, tt.want.link, location)
			}
		})
	}
}

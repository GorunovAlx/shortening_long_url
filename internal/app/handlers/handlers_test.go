package handlers

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStorage imitates ShortURLStorage.
type mockStorage struct {
	id      string
	storage map[string]string
}

func (ms *mockStorage) idInkrement() error {
	id, err := strconv.Atoi(ms.id)
	if err != nil {
		return errors.New(err.Error())
	}
	id += 1
	ms.id = strconv.Itoa(id)
	return nil
}

// Imitating ShortURLRepo.GetInitialLink.
func (ms *mockStorage) GetInitialLink(shortLink string) (string, error) {
	link := ms.storage[shortLink]
	if link == "" {
		return "", errors.New("the url with this value does not exist")
	}
	return link, nil
}

// Imitating ShortURLRepo.CreateShortURL.
func (ms *mockStorage) CreateShortURL(initialLink string) (string, error) {
	ms.storage[ms.id] = initialLink
	defer ms.idInkrement()

	return ms.id, nil
}

// Test request execution.
func testRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func TestCreateShortURLJSONHandler(t *testing.T) {
	target := "http://localhost:8080/"
	type want struct {
		statusCode  int
		contentType string
		expected    string
	}
	tests := []struct {
		name   string
		st     map[string]string
		nextID string
		body   string
		want   want
	}{
		{
			name:   "simple test #1(Post)",
			st:     make(map[string]string),
			nextID: "1",
			body:   `{"url":"google.com"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				expected:    `{"result":"1"}`,
			},
		},
		{
			name: "simple test #2(Post)",
			st: map[string]string{
				"1": "yandex.ru",
				"2": "google.com",
				"3": "tutu.ru",
			},
			nextID: "4",
			body:   `{"url":"yandex.ru"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				expected:    `{"result":"4"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mockStorage{
				id:      tt.nextID,
				storage: tt.st,
			}
			request := httptest.NewRequest(http.MethodPost, target, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateShortURLJSONHandler(&ms))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			buf := new(bytes.Buffer)
			buf.ReadFrom(result.Body)
			str := buf.String()
			err := result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.expected, str)
		})
	}
}

// Test for GetInitialLinkHandler.
func TestGetLinkHandler(t *testing.T) {
	type want struct {
		statusCode int
		link       string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "simple test #1(Get)",
			path: "/2",
			want: want{
				statusCode: 307,
				link:       "google.com",
			},
		},
		{
			name: "simple test #2(Get)",
			path: "/5",
			want: want{
				statusCode: 400,
				link:       "",
			},
		},
		{
			name: "simple test #3(Get)",
			path: "/",
			want: want{
				statusCode: 405,
				link:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mockStorage{
				id: "4",
				storage: map[string]string{
					"1": "yandex.ru",
					"2": "google.com",
					"3": "tutu.ru",
				},
			}
			r := NewHandler(&ms)
			ts := httptest.NewServer(r)
			defer ts.Close()
			result := testRequest(t, ts, "GET", tt.path)

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			location := result.Header.Get("Location")
			assert.Equal(t, tt.want.link, location)
		})
	}
}

// Test for CreateShortURLHandler.
func TestCreateShortURLHandler(t *testing.T) {
	target := "http://localhost:8080/"
	type want struct {
		statusCode int
		link       string
	}
	tests := []struct {
		name   string
		st     map[string]string
		nextID string
		body   []byte
		want   want
	}{
		{
			name:   "simple test #1(Post)",
			st:     make(map[string]string),
			nextID: "1",
			body:   []byte("google.com"),
			want: want{
				statusCode: 201,
				link:       "http://localhost:8080/1",
			},
		},
		{
			name: "simple test #2(Post)",
			st: map[string]string{
				"1": "yandex.ru",
				"2": "google.com",
				"3": "tutu.ru",
			},
			nextID: "4",
			body:   []byte("yandex.ru"),
			want: want{
				statusCode: 201,
				link:       "http://localhost:8080/4",
			},
		},
		{
			name: "simple test #3(Post)",
			st: map[string]string{
				"1": "yandex.ru",
				"2": "google.com",
				"3": "tutu.ru",
			},
			nextID: "4",
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
				id:      tt.nextID,
				storage: tt.st,
			}
			request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateShortURLHandler(&ms))
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			link, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.link, string(link))
		})
	}
}

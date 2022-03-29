package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	gen "github.com/GorunovAlx/shortening_long_url/internal/app/generators"
	mocks "github.com/GorunovAlx/shortening_long_url/internal/app/mocks"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
	"github.com/golang/mock/gomock"
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
func (ms *mockStorage) CreateShortURL(shortURL *storage.ShortURL) (string, error) {
	ms.storage[ms.id] = shortURL.InitialLink
	defer ms.idInkrement()

	return ms.id, nil
}

func (ms *mockStorage) GetAllShortURLUser(id uint32) ([]storage.ShortURLByUser, error) {
	return nil, nil
}

func (ms *mockStorage) PingDB() error {
	return nil
}

func (ms *mockStorage) CreateListShortURL(links []storage.ShortURLByUser) ([]storage.ShortURLByUser, error) {
	return nil, nil
}

// Test request execution.
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, body)
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

// Test for CreateShortURLJSONHandler.
func TestCreateShortURLJSONHandler(t *testing.T) {
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
				expected:    "{\"result\":\"/1\"}",
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
				expected:    "{\"result\":\"/4\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mockStorage{
				id:      tt.nextID,
				storage: tt.st,
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateShortURLJSONHandler(&ms))
			tok, e := gen.GenerateUserIDToken()
			require.NoError(t, e)
			ctx := request.Context()
			ctx = context.WithValue(ctx, contextKeyRequestID, tok)
			h.ServeHTTP(w, request.WithContext(ctx))
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
			r := NewRouter(&ms)
			ts := httptest.NewServer(r)
			defer ts.Close()
			result := testRequest(t, ts, "GET", tt.path, nil)

			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			location := result.Header.Get("Location")
			assert.Equal(t, tt.want.link, location)
		})
	}
}

// Test for CreateShortURLHandler.
func TestCreateShortURLHandler(t *testing.T) {
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
				link:       "/1",
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
				link:       "/4",
			},
		},
		/*
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
					link:       "Incorrect request\n",
				},
			},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//ms := mockStorage{
			//	id:      tt.nextID,
			//	storage: tt.st,
			//}

			var shortURL storage.ShortURL
			shortURL.InitialLink = string(tt.body)

			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockShortURLRepo(ctrl)

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateShortURLHandler(mockStorage))

			token, e := gen.GenerateUserIDToken()
			require.NoError(t, e)
			id, err := gen.GetUserID(token)
			require.NoError(t, err)

			shortURL.UserID = id

			mockStorage.EXPECT().CreateShortURL(&shortURL).Return(tt.want.link[1:], nil)

			ctx := request.Context()
			ctx = context.WithValue(ctx, contextKeyRequestID, token)
			h.ServeHTTP(w, request.WithContext(ctx))
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

func TestGetAllShortURLUserHandler(t *testing.T) {
	var shorts = []storage.ShortURLByUser{
		{
			ShortLink:   "/1",
			InitialLink: "google.com",
		},
		{
			ShortLink:   "/2",
			InitialLink: "yandex.com",
		},
	}

	type want struct {
		statusCode  int
		contentType string
		shorts      []storage.ShortURLByUser
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "simple test #1(Get)",
			want: want{
				statusCode:  200,
				contentType: "application/json",
				shorts:      shorts,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockShortURLRepo(ctrl)

			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetAllShortURLUserHandler(mockStorage))

			token, e := gen.GenerateUserIDToken()
			require.NoError(t, e)
			id, err := gen.GetUserID(token)
			require.NoError(t, err)

			mockStorage.EXPECT().GetAllShortURLUser(id).Return(tt.want.shorts, nil)

			ctx := request.Context()
			ctx = context.WithValue(ctx, contextKeyRequestID, token)
			h.ServeHTTP(w, request.WithContext(ctx))
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			var links []storage.ShortURLByUser
			body, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			json.Unmarshal(body, &links)
			assert.ObjectsAreEqualValues(tt.want.shorts, links)
		})
	}
}

func TestGetPingToDBHandle(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "simple test #1(Get)",
			want: want{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockShortURLRepo(ctrl)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetPingToDBHandle(mockStorage))

			mockStorage.EXPECT().PingDB().Return(nil)

			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestCreateListShortURLHandler(t *testing.T) {
	var bodyTest = `[
		{
			"correlation_id":"d69863e8-958c-4a2b-80bb-e7692779cf0e",
			"original_url":"http://mq8vndl.yandex/jyk7qu4x6jb/v751v/nfh7k5n8dza"
		},
		{
			"correlation_id":"b8be6090-2f46-492a-a98f-1233351a904a",
			"original_url":"http://h4jh4.biz"
		}]`
	var shortsSend = []storage.ShortURLByUser{
		{
			CorrelationID: "d69863e8-958c-4a2b-80bb-e7692779cf0e",
			InitialLink:   "http://mq8vndl.yandex/jyk7qu4x6jb/v751v/nfh7k5n8dza",
		},
		{
			CorrelationID: "b8be6090-2f46-492a-a98f-1233351a904a",
			InitialLink:   "http://h4jh4.biz",
		},
	}
	var shortsGet = []storage.ShortURLByUser{
		{
			CorrelationID: "d69863e8-958c-4a2b-80bb-e7692779cf0e",
			ShortLink:     "/1",
		},
		{
			CorrelationID: "b8be6090-2f46-492a-a98f-1233351a904a",
			ShortLink:     "/2",
		},
	}

	type want struct {
		statusCode  int
		contentType string
		shorts      []storage.ShortURLByUser
	}
	tests := []struct {
		name      string
		shortsFor []storage.ShortURLByUser
		body      string
		want      want
	}{
		{
			name:      "simple test #1(Post)",
			body:      bodyTest,
			shortsFor: shortsSend,
			want: want{
				statusCode:  201,
				contentType: "application/json",
				shorts:      shortsGet,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mocks.NewMockShortURLRepo(ctrl)

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CreateListShortURLHandler(mockStorage))

			mockStorage.EXPECT().CreateListShortURL(tt.shortsFor).Return(tt.want.shorts, nil)

			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			var links []storage.ShortURLByUser
			body, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			json.Unmarshal(body, &links)
			assert.ObjectsAreEqualValues(tt.want.shorts, links)
		})
	}
}

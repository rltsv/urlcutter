package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"github.com/stretchr/testify/assert"
)

func TestHandlerShortener_HeadHandlerPost(t *testing.T) {

	type request struct {
		URL    string
		body   string
		method string
	}
	type want struct {
		body string
		code int
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "test method post with correct initial data",
			request: request{
				URL:  "/api/shorten",
				body: `{"url":"http://postman-echo.com/get"}`,
			},
			want: want{
				code: http.StatusCreated,
				body: `{"result":"http://localhost:8000/1"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			cfg := config.Config{
				ServerAddress:   ":8000",
				BaseURL:         "http://localhost:8000",
				FileStoragePath: "memory.log",
			}

			shortenerRepo := repository.NewStorage(cfg)
			shortenerUsecase := shortener.NewUsecase(*shortenerRepo, cfg)
			handler := NewHandlerShortener(*shortenerUsecase)

			r := httptest.NewRequest(http.MethodPost, tc.request.URL, bytes.NewBufferString(tc.request.body))
			w := httptest.NewRecorder()

			router := SetupRouter(handler)

			router.ServeHTTP(w, r)

			res := w.Result()

			respBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tc.want.code, res.StatusCode)
			assert.Equal(t, tc.want.body, string(respBody))

			os.Args = nil
		})
	}
}

func TestHandlerShortener_HeadHandlerGet(t *testing.T) {
	type request struct {
		shortLink string
	}
	type want struct {
		code         int
		contentField string
	}

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "good initial data",
			request: request{
				shortLink: "/1",
			},
			want: want{
				code:         http.StatusTemporaryRedirect,
				contentField: "http://postman-echo.com/get",
			},
		},
		{
			name: "wrong short link",
			request: request{
				shortLink: "/2",
			},
			want: want{
				code:         http.StatusBadRequest,
				contentField: "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			cfg := config.Config{
				ServerAddress: ":8000",
				BaseURL:       "http://localhost:8000",
			}

			shortenerRepo := repository.Storage{
				InMemoryStorage: map[int]string{
					1: "http://postman-echo.com/get",
				},
				IDCount:   0,
				Mux:       new(sync.RWMutex),
				AppConfig: cfg,
			}
			shortenerUsecase := shortener.NewUsecase(shortenerRepo, cfg)
			handler := NewHandlerShortener(*shortenerUsecase)

			r := httptest.NewRequest(http.MethodGet, tc.request.shortLink, nil)
			w := httptest.NewRecorder()

			router := SetupRouter(handler)
			router.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.want.code, res.StatusCode)
			assert.Equal(t, tc.want.contentField, res.Header.Get("Location"))

		})
	}
}

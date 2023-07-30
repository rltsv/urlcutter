package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	cfg "github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"github.com/stretchr/testify/assert"
)

func TestHandlerShortener_CreateShortLinkViaJSON(t *testing.T) {
	type request struct {
		requestURL    string
		requestBody   string
		requestMethod string
	}
	type config struct {
		serverAddress   string
		baseURL         string
		fileStoragePath string
	}
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		request request
		config  config
		want    want
	}{
		{
			name: "request short link with correct initial data",
			request: request{
				requestURL:  "/api/shorten",
				requestBody: `{"url":"http://postman-echo.com/get"}`,
			},
			config: config{
				serverAddress:   ":8080",
				baseURL:         "http://localhost:8080",
				fileStoragePath: "memory.log",
			},
			want: want{
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			initConfig := cfg.Config{
				ServerAddress:   tc.config.serverAddress,
				BaseURL:         tc.config.baseURL,
				FileStoragePath: tc.config.fileStoragePath,
			}

			var handler *HandlerShortener
			if tc.config.fileStoragePath == "" {
				storage := repository.NewMemoryStorage(initConfig)
				shortenerUsecase := shortener.NewUsecase(storage, initConfig)
				handler = NewHandlerShortener(*shortenerUsecase)
			} else {
				storage := repository.NewFileStorage(initConfig)
				shortenerUsecase := shortener.NewUsecase(storage, initConfig)
				handler = NewHandlerShortener(*shortenerUsecase)
			}

			r := httptest.NewRequest(http.MethodPost, tc.request.requestURL, bytes.NewReader([]byte(tc.request.requestBody)))
			w := httptest.NewRecorder()

			router := SetupRouter(handler)

			router.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.want.statusCode, res.StatusCode)

			os.Args = nil

			os.Remove(tc.config.fileStoragePath)
		})
	}
}

//func TestHandlerShortener_HeadHandlerGet(t *testing.T) {
//	type request struct {
//		shortLink string
//	}
//	type want struct {
//		code         int
//		contentField string
//	}
//
//	tests := []struct {
//		name    string
//		request request
//		want    want
//	}{
//		{
//			name: "good initial data",
//			request: request{
//				shortLink: "/1",
//			},
//			want: want{
//				code:         http.StatusTemporaryRedirect,
//				contentField: "http://postman-echo.com/get",
//			},
//		},
//		{
//			name: "wrong short link",
//			request: request{
//				shortLink: "/2",
//			},
//			want: want{
//				code:         http.StatusBadRequest,
//				contentField: "",
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//
//			cfg := config.Config{
//				ServerAddress: ":8000",
//				BaseURL:       "http://localhost:8000",
//			}
//			1: "http://postman-echo.com/get"
//			shortenerRepo := repository.MemoryStorage{
//				Links: []entity.Link{
//
//				},
//				Mux:       new(sync.RWMutex),
//				AppConfig: cfg,
//			}
//			shortenerUsecase := shortener.NewUsecase(shortenerRepo, cfg)
//			handler := NewHandlerShortener(*shortenerUsecase)
//
//			r := httptest.NewRequest(http.MethodGet, tc.request.shortLink, nil)
//			w := httptest.NewRecorder()
//
//			router := SetupRouter(handler)
//			router.ServeHTTP(w, r)
//
//			res := w.Result()
//			defer res.Body.Close()
//
//			assert.Equal(t, tc.want.code, res.StatusCode)
//			assert.Equal(t, tc.want.contentField, res.Header.Get("Location"))
//
//		})
//	}
//}

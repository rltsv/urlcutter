package rest

import (
	"bytes"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestHandlerShortener_HeadHandler_MethodPost(t *testing.T) {

	type request struct {
		URL  string
		body string
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
				body: `{"result":"http://localhost:8080/1"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			shortenerRepo := repository.NewLinksRepository()
			shortenerUsecase := shortener.NewUsecase(shortenerRepo)
			handler := NewHandlerShortener(*shortenerUsecase)

			err := config.InitConfig()
			assert.NoError(t, err)

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

		})
	}

}

func TestHandlerShortener_HeadHandler_MethodGet(t *testing.T) {
	type request struct {
		body      string
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
			name: "test method get",
			request: request{
				body:      "http://postman-echo.com/get",
				shortLink: "/1",
			},
			want: want{
				code:         http.StatusTemporaryRedirect,
				contentField: "http://postman-echo.com/get",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shortenerRepo := repository.LinksRepository{
				Storage: map[int]string{
					1: "http://postman-echo.com/get",
				},
				IDCount: 1,
				Mux:     new(sync.RWMutex),
			}
			shortenerUsecase := shortener.NewUsecase(&shortenerRepo)
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

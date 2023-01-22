package rest

import (
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestHandlerShortener_HeadHandler_MethodPost(t *testing.T) {
	type request struct {
		URL    string
		patURL string
		body   string
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
				URL:  "http://localhost:8080/",
				body: "http://postman-echo.com/get",
			},
			want: want{
				code: 201,
				body: "http://localhost:8080/1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			shortenerRepo := repository.NewLinksRepository()
			shortenerUsecase := shortener.NewUsecase(shortenerRepo)
			handler := NewHandlerShortener(*shortenerUsecase)

			myReader := strings.NewReader(tc.request.body)

			r := httptest.NewRequest(http.MethodPost, tc.request.URL, myReader)
			w := httptest.NewRecorder()

			router := SetupRouter(handler)

			router.ServeHTTP(w, r)

			res := w.Result()
			res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.want.code, res.StatusCode)
			assert.Equal(t, tc.want.body, string(resBody))

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
				Mux:     &sync.Mutex{},
			}
			shortenerUsecase := shortener.NewUsecase(&shortenerRepo)
			handler := NewHandlerShortener(*shortenerUsecase)

			r := httptest.NewRequest(http.MethodGet, tc.request.shortLink, nil)
			w := httptest.NewRecorder()

			router := SetupRouter(handler)
			router.ServeHTTP(w, r)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, tc.want.code, res.StatusCode)
			assert.Equal(t, tc.want.contentField, res.Header.Get("Location"))

		})
	}
}

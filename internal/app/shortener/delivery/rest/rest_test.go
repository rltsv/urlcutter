package rest

import (
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
			shortenerHandler := NewHandlerShortener(*shortenerUsecase)

			myReader := strings.NewReader(tc.request.body)

			r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", myReader)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(shortenerHandler.HeadHandler)
			h.ServeHTTP(w, r)

			res := w.Result()

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
		mainPageURL  string
		body         string
		shortLinkURL string
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
			name: "test method post with correct initial data",
			request: request{
				mainPageURL:  "http://localhost:8080/",
				body:         "http://postman-echo.com/get",
				shortLinkURL: "http://localhost:8080/1",
			},
			want: want{
				code:         http.StatusTemporaryRedirect,
				contentField: "http://postman-echo.com/get",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shortenerRepo := repository.NewLinksRepository()
			shortenerUsecase := shortener.NewUsecase(shortenerRepo)
			shortenerHandler := NewHandlerShortener(*shortenerUsecase)

			myReader := strings.NewReader(tc.request.body)

			requestPost := httptest.NewRequest(http.MethodPost, tc.request.mainPageURL, myReader)
			requestGet := httptest.NewRequest(http.MethodGet, tc.request.shortLinkURL, nil)
			recPost := httptest.NewRecorder()
			recGet := httptest.NewRecorder()

			h := http.HandlerFunc(shortenerHandler.HeadHandler)
			h.ServeHTTP(recPost, requestPost)
			h.ServeHTTP(recGet, requestGet)

			result := recGet.Result()

			assert.Equal(t, tc.want.code, result.StatusCode)

		})
	}
}

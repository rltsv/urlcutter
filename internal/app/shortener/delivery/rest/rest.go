package rest

import (
	"context"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HandlerShortener struct {
	useCase shortener.Usecase
}

func NewHandlerShortener(useCase shortener.Usecase) *HandlerShortener {
	return &HandlerShortener{
		useCase: useCase,
	}
}

func (hs *HandlerShortener) HeadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method == http.MethodPost {
		if strings.TrimLeft(r.URL.Path, "/") != "" {
			http.Error(w, "specify the request", 400)
			return
		}

		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "specify the request", 400)
			return
		}

		if len(respBody) == 0 {
			http.Error(w, "where is nothing to short, check body", 400)
			return
		}

		shortLink := hs.useCase.CreateShortLink(ctx, string(respBody))

		w.WriteHeader(201)
		_, err = w.Write([]byte(shortLink))
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Fatal("", err)
			return
		}

	} else if r.Method == http.MethodGet {

		idValue := strings.TrimPrefix(r.URL.Path, "/")

		if idValue == "" {
			http.Error(w, "specify the request", 400)
			return
		}
		id, err := strconv.Atoi(idValue)
		if err != nil {
			http.Error(w, "specify the request", 400)
			return
		}

		origLink, err := hs.useCase.GetLinkByID(ctx, id)
		if err == repository.ErrLinkNotFound {
			http.Error(w, "there is no any link by this id", 400)
		}

		w.Header().Set("Location", origLink)
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else {
		http.Error(w, "request error", 400)
	}

}

package main

import (
	"github.com/rltsv/urlcutter/internal/app/shortener/delivery/rest"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"log"
	"net/http"
)

func main() {
	repo := repository.NewLinksRepository()
	repoUsecase := shortener.NewUsecase(repo)
	handler := rest.NewHandlerShortener(*repoUsecase)

	router := rest.SetupRouter(handler)

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}

package main

import (
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/delivery/rest"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"log"
	"net/http"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewLinksRepository()
	repoUsecase := shortener.NewUsecase(repo)
	handler := rest.NewHandlerShortener(*repoUsecase)

	router := rest.SetupRouter(handler)

	log.Printf("app starts listen on port: %s", config.Cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(":"+config.Cfg.ServerAddress, router))
}

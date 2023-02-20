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

	shortenerStorage := repository.NewStorage()
	shortenerUsecase := shortener.NewUsecase(*shortenerStorage)
	handler := rest.NewHandlerShortener(*shortenerUsecase)

	router := rest.SetupRouter(handler)

	log.Printf("app starts listen on : %s", config.Cfg.ServerAddress)
	log.Printf("BASE_URL is : %s", config.Cfg.BaseURL)
	log.Fatal(http.ListenAndServe(config.Cfg.ServerAddress, router))
}

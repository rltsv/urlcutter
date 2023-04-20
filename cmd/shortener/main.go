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

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	memoryStorage := repository.NewMemoryStorage(cfg)
	fileStorage := repository.NewFileStorage(cfg)
	shortenerUsecase := shortener.NewUsecase(memoryStorage, fileStorage, cfg)
	handler := rest.NewHandlerShortener(*shortenerUsecase)

	router := rest.SetupRouter(handler)

	log.Printf("http server startup address is %s", cfg.ServerAddress)
	log.Printf("the base address of the resulting shortened URL : %s", cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}

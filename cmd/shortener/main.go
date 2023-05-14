package main

import (
	"log"
	"net/http"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/delivery/rest"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	var dbStorage *repository.PsqlStorage
	if cfg.DataBaseDSN != "" {
		db, err := config.InitDB(cfg)
		if err != nil {
			log.Fatal(err)
		}
		dbStorage = repository.NewPsqlStorage(db)
	}

	var handler *rest.HandlerShortener
	if cfg.FileStoragePath == "" {
		storage := repository.NewMemoryStorage(cfg)
		shortenerUsecase := shortener.NewUsecase(storage, dbStorage, cfg)
		handler = rest.NewHandlerShortener(*shortenerUsecase)
	} else {
		storage := repository.NewFileStorage(cfg)
		shortenerUsecase := shortener.NewUsecase(storage, dbStorage, cfg)
		handler = rest.NewHandlerShortener(*shortenerUsecase)
	}

	router := rest.SetupRouter(handler)

	log.Printf("http server startup address is %s", cfg.ServerAddress)
	log.Printf("the base address of the resulting shortened URL : %s", cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}

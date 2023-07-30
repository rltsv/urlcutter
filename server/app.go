package server

import (
	"log"

	"github.com/gin-gonic/gin"
	cfg "github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/delivery/rest"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
)

func InitApp(config cfg.Config) (*gin.Engine, error) {

	var handler *rest.HandlerShortener
	if config.DataBaseDSN != "" {
		db, err := cfg.InitDB(config)
		if err != nil {
			log.Fatal(err)
		}
		dbStorage := repository.NewPsqlStorage(db)
		shortenerService := shortener.NewUsecase(dbStorage, config)
		handler = rest.NewHandlerShortener(*shortenerService)
	} else if config.FileStoragePath != "" {
		storage := repository.NewFileStorage(config)
		shortenerService := shortener.NewUsecase(storage, config)
		handler = rest.NewHandlerShortener(*shortenerService)
	} else {
		storage := repository.NewMemoryStorage(config)
		shortenerService := shortener.NewUsecase(storage, config)
		handler = rest.NewHandlerShortener(*shortenerService)
	}

	return rest.SetupRouter(handler), nil
}

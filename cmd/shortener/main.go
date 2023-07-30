package main

import (
	"log"
	"net/http"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/server"
)

func main() {

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	router, err := server.InitApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("http server startup address is %s", cfg.ServerAddress)
	log.Printf("the base address of the resulting shortened URL : %s", cfg.BaseURL)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}

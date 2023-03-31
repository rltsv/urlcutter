package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func InitConfig() (cfg Config, err error) {
	config := Config{}

	err = env.Parse(&config)
	if err != nil {
		return Config{}, err
	}

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); !ok {
		flag.StringVar(&config.ServerAddress, "a", ":8080", "http server startup address")
	} else {
		log.Printf("environment variable SERVER_ADDRESS is set as - %s", val)
	}

	if val, ok := os.LookupEnv("BASE_URL"); !ok {
		flag.StringVar(&config.BaseURL, "b", config.ServerAddress, "the base address of the resulting shortened URL")
	} else {
		log.Printf("environment variable BASE_URL is set as - %s", val)
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); !ok {
		flag.StringVar(&config.FileStoragePath, "f", "", "the path to the file with the abbreviated URL")
	} else {
		log.Printf("environment variable SERVER_ADDRESS is set as - %s", val)
	}

	flag.Parse()

	return config, nil
}

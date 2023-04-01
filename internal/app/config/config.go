package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func InitConfig() (cfg Config, err error) {
	config := Config{}

	srvrAddr := flag.String("a", "", "http server startup address")
	baseURL := flag.String("b", "", "the base address of the resulting shortened URL")
	filePath := flag.String("f", "", "the path to the file with the abbreviated URL")
	flag.Parse()

	err = env.Parse(&config)
	if err != nil {
		return Config{}, err
	}

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); !ok {
		if *srvrAddr != "" {
			config.ServerAddress = *srvrAddr
		}
	} else {
		log.Printf("environment variable SERVER_ADDRESS is set as - %s", val)

	}

	if val, ok := os.LookupEnv("BASE_URL"); !ok {
		if *baseURL != "" {
			config.BaseURL = *baseURL
		} else if _, ok = os.LookupEnv("SERVER_ADDRESS"); ok {
			config.BaseURL = "http://" + config.ServerAddress
		} else if _, ok = os.LookupEnv("SERVER_ADDRESS"); !ok {
			if config.ServerAddress == *srvrAddr {
				config.BaseURL = "http://" + *srvrAddr
			}
		}
	} else {
		log.Printf("environment variable BASE_URL is set as - %s", val)
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); !ok {
		if *filePath != "" {
			config.FileStoragePath = *filePath
		}
	} else {
		log.Printf("environment variable FILE_STORAGE_PATH is set as - %s", val)
	}

	return config, nil
}

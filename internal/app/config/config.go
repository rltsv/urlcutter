package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS,required"`
	BaseURL       string `env:"BASE_URL,required"`
}

var Cfg Config

func InitConfig() error {
	err := env.Parse(&Cfg)
	if err != nil {
		return err
	}
	return nil
}

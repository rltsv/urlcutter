package config

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
}

func InitConfig() (cfg Config, err error) {
	config := Config{}

	// setup flags for later usage
	srvrAddr := flag.String("a", "", "http server startup address")
	baseURL := flag.String("b", "", "the base address of the resulting shortened URL")
	filePath := flag.String("f", "", "the path to the file with the abbreviated URL")
	db := flag.String("d", "", "address for connection to db")
	flag.Parse()

	err = env.Parse(&config)
	if err != nil {
		return Config{}, err
	}
	// check server address config
	if val, ok := os.LookupEnv("SERVER_ADDRESS"); !ok {
		if *srvrAddr != "" {
			config.ServerAddress = *srvrAddr
		}
	} else {
		log.Printf("environment variable SERVER_ADDRESS is set as - %s", val)

	}
	// check base url config
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
	// check should we use file for storage or not
	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); !ok {
		if *filePath != "" {
			config.FileStoragePath = *filePath
		}
	} else {
		log.Printf("environment variable FILE_STORAGE_PATH is set as - %s", val)
	}

	if val, ok := os.LookupEnv("DATABASE_DSN"); !ok {
		if *db != "" {
			config.DataBaseDSN = *db
		}
	} else {
		log.Printf("environment variable DATABASE_DSN is set as - %s", val)
	}
	return config, nil
}

func InitDB(cfg Config) (conn *pgx.Conn, err error) {

	conn, err = pgx.Connect(context.Background(), cfg.DataBaseDSN)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	var connection = "ok"
	err = conn.Ping(ctx)
	if err != nil {
		connection = "false"
	}
	log.Printf("connection to db is %s", connection)

	file, err := os.OpenFile("db.sql", os.O_RDONLY, 0777)
	if err != nil {
		log.Println(err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}

	_, err = conn.Exec(ctx, string(bytes))
	if err != nil {
		log.Print(err)
	}

	return conn, nil
}

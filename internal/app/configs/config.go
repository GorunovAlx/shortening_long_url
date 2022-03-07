package configs

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

var Cfg Config

func SetConfigs() {
	parameters := os.Args[1:]
	log.Println(parameters)

	err := env.Parse(&Cfg)
	if err != nil {
		log.Println(err)
	}

	if len(parameters) > 0 {
		if flag.Lookup("a") == nil {
			flag.StringVar(&Cfg.ServerAddress, "a", Cfg.ServerAddress, "server address to listen on")
		}
		if flag.Lookup("b") == nil {
			flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "base url to listen on")
		}
		if flag.Lookup("f") == nil {
			flag.StringVar(&Cfg.FileStoragePath, "f", Cfg.FileStoragePath, "file storage path")
		}
		flag.Parse()
	}

	log.Println(Cfg)
}

package configs

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

var Cfg Config

func SetConfigs() error {
	parameters := os.Args[1:]
	log.Println(parameters)
	if len(parameters) > 0 {
		pflag.StringVarP(&Cfg.ServerAddress, "a", "a", ":8080", "server address to listen on")
		pflag.StringVarP(&Cfg.BaseURL, "b", "b", "http://localhost:8080", "base url to listen on")
		pflag.StringVarP(&Cfg.FileStoragePath, "f", "f", "", "file storage path")
		pflag.Parse()
		log.Println(Cfg)
		return nil
	} else {
		err := env.Parse(&Cfg)
		if err != nil {
			return err
		}
		log.Println(Cfg)
	}
	log.Println(Cfg)
	return nil
}

package configs

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://${SERVER_ADDRESS}"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

var Cfg Config

func SetConfigs() error {
	parameters := os.Args[1:]
	log.Println(parameters)
	if len(parameters) > 0 {
		pflag.StringVarP(&Cfg.ServerAddress, "a", "a", Cfg.ServerAddress, "server address to listen on")
		pflag.StringVarP(&Cfg.BaseURL, "b", "b", Cfg.BaseURL, "base url to listen on")
		pflag.StringVarP(&Cfg.FileStoragePath, "f", "f", Cfg.FileStoragePath, "file storage path")
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

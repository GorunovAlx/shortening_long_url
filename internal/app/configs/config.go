package configs

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

// The config type is a structure containing:
// ServerAddress - the server address,
// BaseURL - the base address of the resulting shortened url,
// FileStoragePath - the path to the file where the shortened url is written.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:""`
	SecretKey       string `env:"SECRET_KEY" envDefault:"secret_key"`
}

var Cfg Config

// Parsing the environment variables and returns an error, interrupting program execution.
// Checks if flags are passed, the config will be initialized.
func SetConfig() {
	parameters := os.Args[1:]

	if err := env.Parse(&Cfg); err != nil {
		log.Fatal(err)
	}

	if len(parameters) > 0 {
		if flag.Lookup("a") == nil {
			flag.StringVar(&Cfg.ServerAddress, "a", Cfg.ServerAddress, "HTTP server launch address")
		}
		if flag.Lookup("b") == nil {
			flag.StringVar(&Cfg.BaseURL, "b", Cfg.BaseURL, "the base address of the resulting shortened URL")
		}
		if flag.Lookup("f") == nil {
			flag.StringVar(&Cfg.FileStoragePath, "f", Cfg.FileStoragePath, "file storage path")
		}
		if flag.Lookup("d") == nil {
			flag.StringVar(&Cfg.DatabaseDSN, "d", Cfg.DatabaseDSN, "string with connection address to db")
		}
		flag.Parse()
	}
}

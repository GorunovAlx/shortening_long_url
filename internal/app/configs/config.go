package configs

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	// The server address
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	// The base address of the resulting shortened url
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	// The path to the file where the shortened url is written.
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	// Database url string like postgres://postgres:pass@localhost:5432/dbname?sslmode=disable
	DatabaseDSN string `env:"DATABASE_DSN" envDefault:""`
	// Secret key to encrypt data
	SecretKey string `env:"SECRET_KEY" envDefault:"secret_key"`
	// logging level for zerolog
	ZerologLevel int8 `env:"ZERO_LOG_LEVEL" envDefault:"0"`
	// logging level for pgx driver db
	PgxLogLevel string `env:"PGX_LOG_LEVEL" envDefualt:"info"`
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

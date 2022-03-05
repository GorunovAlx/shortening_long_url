package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func Init() {
	flag.StringVar(&handlers.Cfg.ServerAddress, "a", "localhost:8080", "server address to listen on")
	flag.Lookup("a").NoOptDefVal = handlers.Cfg.ServerAddress
	flag.StringVar(&handlers.Cfg.BaseURL, "b", "http://localhost:8080", "base url to listen on")
	flag.Lookup("b").NoOptDefVal = handlers.Cfg.BaseURL
	flag.StringVar(&storage.Cfg.FileStoragePath, "f", "", "file storage path")
	flag.Lookup("c").NoOptDefVal = storage.Cfg.FileStoragePath
}

func main() {
	Init()
	flag.Parse()
	log.Println(os.Args)
	urlStorage, err := storage.NewShortURLStorage()
	if err != nil {
		log.Fatal(errors.New("an error occurred while creating the repository "))
	}
	handler := handlers.NewHandler(urlStorage)
	log.Println(os.LookupEnv("FILE_STORAGE_PATH"))
	log.Fatal(http.ListenAndServe(handlers.Cfg.ServerAddress, handler))
}

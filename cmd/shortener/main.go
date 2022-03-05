package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func Init() {
	flag.StringVar(&handlers.Cfg.ServerAddress, "a", handlers.Cfg.ServerAddress, "server address to listen on")
	flag.StringVar(&handlers.Cfg.BaseURL, "b", handlers.Cfg.BaseURL, "base url to listen on")
	flag.StringVar(&storage.Cfg.FileStoragePath, "f", storage.Cfg.FileStoragePath, "file storage path")
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

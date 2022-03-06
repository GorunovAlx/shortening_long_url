package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func main() {
	configs.SetConfigs()
	urlStorage, err := storage.NewShortURLStorage()
	if err != nil {
		log.Fatal(errors.New("an error occurred while creating the repository "))
	}
	handler := handlers.NewHandler(urlStorage)
	log.Fatal(http.ListenAndServe(configs.Cfg.ServerAddress, handler))
}

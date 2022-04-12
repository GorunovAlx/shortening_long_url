package main

import (
	"log"
	"net/http"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
	"github.com/GorunovAlx/shortening_long_url/internal/app/utils"
)

func main() {
	configs.SetConfig()
	utils.LoggerInit()
	urlStorage := storage.NewStorage()
	handler := handlers.NewRouter(urlStorage)
	log.Fatal(http.ListenAndServe(configs.Cfg.ServerAddress, handler))
}

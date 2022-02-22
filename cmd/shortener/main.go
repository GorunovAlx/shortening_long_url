package main

import (
	"log"

	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func main() {
	urlStorage := storage.NewShortURLStorage()
	r := handlers.RegisterRoutes(urlStorage)
	log.Fatal(r.Run())
}

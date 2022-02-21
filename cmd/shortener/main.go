package main

import (
	"log"
	"net/http"

	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func main() {
	urlStorage := storage.NewShortURLStorage()
	//handler := handlers.NewHandler(urlStorage)
	r := handlers.RegisterRoutes(urlStorage)
	log.Fatal(http.ListenAndServe(":8080", r))
}

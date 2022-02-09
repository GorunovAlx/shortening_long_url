package main

import (
	"log"
	"net/http"

	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func main() {
	urlStorage := storage.NewShortURLStorage()
	http.HandleFunc("/", handlers.ShortenerURLHandler(urlStorage))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

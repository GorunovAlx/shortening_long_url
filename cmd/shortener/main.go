package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func main() {
	urlStorage := storage.NewShortURLStorage()
	handler := handlers.NewHandler(urlStorage)
	adr := os.Getenv("SERVER_ADDRESS")
	log.Fatal(http.ListenAndServe(adr, handler))
}

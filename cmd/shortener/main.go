package main

import (
	"log"
	"net/http"

	"github.com/GorunovAlx/shortening_long_url/internal/app/routers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
)

func main() {
	urlStorage := storage.NewShortURLStorage()
	handler := routers.NewHandler(urlStorage)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

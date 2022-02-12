package routers

import (
	"github.com/GorunovAlx/shortening_long_url/internal/app/handlers"
	"github.com/GorunovAlx/shortening_long_url/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	*chi.Mux
	Repo storage.ShortURLRepo
}

func NewHandler(repo storage.ShortURLRepo) *Handler {
	h := &Handler{
		Mux:  chi.NewMux(),
		Repo: repo,
	}
	h.Post("/", handlers.CreateShortURLHandler(repo))
	h.Get("/{shortURL}", handlers.GetInitialLinkHandler(repo))

	return h
}

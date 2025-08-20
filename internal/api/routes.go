package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Post("/currency/add", h.AddCurrency)
	r.Post("/currency/remove", h.RemoveCurrency)
	r.Get("/currency/price", h.GetPrice)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}

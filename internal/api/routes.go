package api

import (
	"net/http"

	"crypto-observer/internal/service"

	_ "crypto-observer/docs"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(s service.CurrencyService) http.Handler {
	h := NewHandler(s)
	r := chi.NewRouter()
	r.Post("/currency/add", h.AddCurrency)
	r.Post("/currency/remove", h.RemoveCurrency)
	r.Get("/currency/price", h.GetPrice)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}

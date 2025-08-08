package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"crypto-observer/internal/logger"
	"crypto-observer/internal/model"
	"crypto-observer/internal/service"
)

type Handler struct {
	service service.CurrencyService
}

func NewHandler(s service.CurrencyService) *Handler { return &Handler{service: s} }

// POST /currency/add
func (h *Handler) AddCurrency(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Symbol string `json:"symbol"`
		Period int    `json:"period"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Warn("AddCurrency: bad request")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	logger.Log.WithFields(logger.Fields{"symbol": req.Symbol, "period": req.Period}).Info("AddCurrency: start")

	if err := h.service.AddCurrency(req.Symbol, req.Period); err != nil {
		logger.Log.WithError(err).Error("AddCurrency: service failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// POST /currency/remove
func (h *Handler) RemoveCurrency(w http.ResponseWriter, r *http.Request) {
	var req model.RemoveCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Warn("RemoveCurrency: bad request")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	logger.Log.WithField("symbol", req.Symbol).Info("RemoveCurrency: start")

	if err := h.service.RemoveCurrency(req.Symbol); err != nil {
		logger.Log.WithError(err).Error("RemoveCurrency: service failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

var (
	ErrInvalidTimestamp = errors.New("invalid timestamp")
	ErrPriceNotFound    = errors.New("price not found")
)

func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	ts := r.URL.Query().Get("timestamp")

	logger.Log.WithFields(logger.Fields{"symbol": symbol, "timestamp": ts}).Info("GetPrice: start")

	resp, err := h.service.GetPrice(symbol, ts)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTimestamp):
			http.Error(w, "invalid timestamp", http.StatusBadRequest) // 400
		case errors.Is(err, service.ErrPriceNotFound):
			http.Error(w, "not found", http.StatusNotFound) // 404
		default:
			logger.Log.WithError(err).Error("unexpected error in GetPrice")
			http.Error(w, "internal error", http.StatusInternalServerError) // 500
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

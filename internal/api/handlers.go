package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"crypto-observer/internal/model"
	"crypto-observer/pkg/logger"
)

type Handler struct {
	service CurrencyService
}

func NewHandler(s CurrencyService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) AddCurrency(w http.ResponseWriter, r *http.Request) {
	var req model.AddReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.L().WithError(err).Warn("AddCurrency: bad request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}
	if err := h.service.AddCurrency(req.Symbol, req.Period); err != nil {
		logger.L().WithError(err).Error("AddCurrency: service failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.StatusOK())
}

func (h *Handler) RemoveCurrency(w http.ResponseWriter, r *http.Request) {
	var req model.RemoveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.L().WithError(err).Warn("RemoveCurrency: bad request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}
	if err := h.service.RemoveCurrency(req.Symbol); err != nil {
		logger.L().WithError(err).Error("RemoveCurrency: service failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	tsStr := r.URL.Query().Get("timestamp")

	logger.L().WithFields(logger.Fields{
		"symbol": symbol,
		"ts":     tsStr,
	}).Info("GetPrice: start")

	if symbol == "" {
		http.Error(w, "symbol is required", http.StatusBadRequest)
		return
	}

	var ts int64
	if tsStr != "" {
		v, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid timestamp", http.StatusBadRequest)
			return
		}
		ts = v
	}

	price, err := h.service.GetPrice(symbol, ts)
	if err != nil {
		// ЛЮБАЯ внутренняя ошибка сервиса → 500
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if price == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := model.PriceDTO{
		Coin:      price.Symbol,
		Timestamp: price.TS,
		Price:     price.Price,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

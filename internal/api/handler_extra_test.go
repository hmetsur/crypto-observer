package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"unicode"

	"crypto-observer/internal/model"
	"crypto-observer/internal/service"

	"github.com/stretchr/testify/require"
)

// Мок CurrencyService, который валидирует вход
type dummyService struct{}

func (d *dummyService) AddCurrency(symbol string, period int) error { return nil }
func (d *dummyService) RemoveCurrency(symbol string) error          { return nil }
func (d *dummyService) GetPrice(symbol, ts string) (*model.PriceResponse, error) {
	// Если timestamp пустой или не число — возвращаем ErrInvalidTimestamp -> 400
	if ts == "" || !isDigits(ts) {
		return nil, service.ErrInvalidTimestamp
	}
	// В остальных случаях считаем, что всё ок
	return &model.PriceResponse{Coin: symbol, Timestamp: 1, Price: 1}, nil
}

func isDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

func TestGetPrice_MissingParams(t *testing.T) {
	h := NewHandler(&dummyService{})

	req := httptest.NewRequest(http.MethodGet, "/currency/price", nil) // нет symbol и timestamp
	w := httptest.NewRecorder()

	h.GetPrice(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code) // 400 из-за ErrInvalidTimestamp
}

func TestGetPrice_BadTimestamp(t *testing.T) {
	h := NewHandler(&dummyService{})

	req := httptest.NewRequest(http.MethodGet, "/currency/price?symbol=btc&timestamp=abc", nil)
	w := httptest.NewRecorder()

	h.GetPrice(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code) // 400 из-за ErrInvalidTimestamp
}

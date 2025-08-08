package api

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"crypto-observer/internal/model"
	"crypto-observer/internal/service"
)

// сервис, который всегда возвращает ошибки — чтобы проверить 400/500/404
type badService struct{}

func (b *badService) AddCurrency(string, int) error { return errors.New("boom") }
func (b *badService) RemoveCurrency(string) error   { return errors.New("boom") }
func (b *badService) GetPrice(string, string) (*model.PriceResponse, error) {
	return nil, service.ErrPriceNotFound
}

var _ service.CurrencyService = (*badService)(nil)

func TestAddCurrency_BadJSON(t *testing.T) {
	h := NewHandler(&badService{})
	req := httptest.NewRequest("POST", "/currency/add", bytes.NewBufferString("{"))
	rr := httptest.NewRecorder()

	h.AddCurrency(rr, req)
	if rr.Code != 400 {
		t.Fatalf("ожидался 400, получили %d", rr.Code)
	}
}

func TestAddCurrency_ServiceError(t *testing.T) {
	h := NewHandler(&badService{})
	body := bytes.NewBufferString(`{"symbol":"btc","period":5}`)
	req := httptest.NewRequest("POST", "/currency/add", body)
	rr := httptest.NewRecorder()

	h.AddCurrency(rr, req)
	if rr.Code != 500 {
		t.Fatalf("ожидался 500, получили %d", rr.Code)
	}
}

func TestRemoveCurrency_ServiceError(t *testing.T) {
	h := NewHandler(&badService{})
	body := bytes.NewBufferString(`{"symbol":"btc"}`)
	req := httptest.NewRequest("POST", "/currency/remove", body)
	rr := httptest.NewRecorder()

	h.RemoveCurrency(rr, req)
	if rr.Code != 500 {
		t.Fatalf("ожидался 500, получили %d", rr.Code)
	}
}

func TestGetPrice_NotFound(t *testing.T) {
	h := NewHandler(&badService{})
	req := httptest.NewRequest("GET", "/currency/price?symbol=btc&timestamp=1", nil)
	rr := httptest.NewRecorder()

	h.GetPrice(rr, req)
	if rr.Code != 404 {
		t.Fatalf("ожидался 404, получили %d", rr.Code)
	}
}

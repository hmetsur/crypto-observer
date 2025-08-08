package api

import (
	"bytes"
	"crypto-observer/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeService struct{}

func (f *fakeService) AddCurrency(symbol string, period int) error { return nil }
func (f *fakeService) RemoveCurrency(symbol string) error          { return nil }
func (f *fakeService) GetPrice(symbol, ts string) (*model.PriceResponse, error) {
	return &model.PriceResponse{Coin: symbol, Timestamp: 1, Price: 100}, nil
}

func TestHandler_AddCurrency(t *testing.T) {
	h := NewHandler(&fakeService{})
	req := httptest.NewRequest("POST", "/currency/add", bytes.NewReader([]byte(`{"symbol":"btc","period":5}`)))
	w := httptest.NewRecorder()
	h.AddCurrency(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("AddCurrency: ожидался 200, а был %d", w.Result().StatusCode)
	}
}

func TestHandler_RemoveCurrency(t *testing.T) {
	h := NewHandler(&fakeService{})
	req := httptest.NewRequest("POST", "/currency/remove", bytes.NewReader([]byte(`{"symbol":"btc"}`)))
	w := httptest.NewRecorder()
	h.RemoveCurrency(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("RemoveCurrency: ожидался 200, а был %d", w.Result().StatusCode)
	}
}

func TestHandler_GetPrice(t *testing.T) {
	h := NewHandler(&fakeService{})
	req := httptest.NewRequest("GET", "/currency/price?symbol=btc&timestamp=1", nil)
	w := httptest.NewRecorder()
	h.GetPrice(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("GetPrice: ожидался 200, а был %d", w.Result().StatusCode)
	}
	var resp model.PriceResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Coin != "btc" {
		t.Errorf("GetPrice: ожидал btc, получил %s", resp.Coin)
	}
}

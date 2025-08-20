package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-observer/internal/model"
)

type fakeServ struct {
	addSymbol    string
	addPeriod    int
	removeSymbol string

	priceResp *model.Price
	priceErr  error
}

func (f *fakeServ) AddCurrency(symbol string, periodSec int) error {
	f.addSymbol = symbol
	f.addPeriod = periodSec
	return nil
}

func (f *fakeServ) RemoveCurrency(symbol string) error {
	f.removeSymbol = symbol
	return nil
}

func (f *fakeServ) GetPrice(symbol string, ts int64) (*model.Price, error) {
	return f.priceResp, f.priceErr
}

// -------------------------------------------------------------

func TestNewRouter_AddCurrency(t *testing.T) {
	svc := &fakeServ{}
	h := NewHandler(svc)
	r := NewRouter(h)

	body := []byte(`{"symbol":"btc","period":5}`)
	req := httptest.NewRequest(http.MethodPost, "/currency/add", bytes.NewReader(body)).
		WithContext(context.Background())
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d; body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
	if svc.addSymbol != "btc" || svc.addPeriod != 5 {
		t.Fatalf("service.AddCurrency called with: %q,%d; want %q,%d",
			svc.addSymbol, svc.addPeriod, "btc", 5)
	}
}

func TestNewRouter_RemoveCurrency(t *testing.T) {
	svc := &fakeServ{}
	h := NewHandler(svc)
	r := NewRouter(h)

	body := []byte(`{"symbol":"eth"}`)
	req := httptest.NewRequest(http.MethodPost, "/currency/remove", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d; body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}
	if svc.removeSymbol != "eth" {
		t.Fatalf("service.RemoveCurrency called with %q; want %q", svc.removeSymbol, "eth")
	}
}

func TestNewRouter_GetPrice(t *testing.T) {
	want := &model.Price{Symbol: "btc", TS: 111, Price: 12345}
	svc := &fakeServ{priceResp: want}
	h := NewHandler(svc)
	r := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/currency/price?symbol=btc&timestamp=111", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d; body=%s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var got model.PriceDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Coin != want.Symbol || got.Timestamp != want.TS || got.Price != want.Price {
		t.Fatalf("response mismatch: %#v vs %#v", got, want)
	}
}

func TestNewRouter_SwaggerMounted(t *testing.T) {
	svc := &fakeServ{}
	h := NewHandler(svc)
	r := NewRouter(h)

	// смоук: маршрут смонтирован (не 404)
	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	if rr.Code == http.StatusNotFound {
		t.Fatalf("swagger route seems not mounted, got 404")
	}
}

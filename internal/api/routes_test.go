package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-observer/internal/model"
	"crypto-observer/internal/service"
)

type emptyService struct{}

func (*emptyService) AddCurrency(string, int) error { return nil }
func (*emptyService) RemoveCurrency(string) error   { return nil }
func (*emptyService) GetPrice(string, string) (*model.PriceResponse, error) {
	return &model.PriceResponse{}, nil
}

var _ service.CurrencyService = (*emptyService)(nil)

func TestNewRouter_RoutesExist(t *testing.T) {
	r := NewRouter(&emptyService{})

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/currency/add"},
		{"POST", "/currency/remove"},
		{"GET", "/currency/price"},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if rr.Code == http.StatusNotFound {
			t.Fatalf("маршрут %s %s не зарегистрирован", tc.method, tc.path)
		}
	}
}

package coingecko

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPriceUSD_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"btc":{"usd":123.45}}`))
	}))
	defer srv.Close()

	c := New(srv.URL, time.Second)
	cents, err := c.GetPriceCents(context.Background(), "btc")
	if err != nil {
		t.Fatalf("GetPriceUSD: %v", err)
	}
	if cents != 12345 {
		t.Fatalf("want 12345, got %d", cents)
	}
}

func TestGetPriceUSD_BadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := New(srv.URL, time.Second)
	if _, err := c.GetPriceCents(context.Background(), "btc"); err == nil {
		t.Fatalf("expected error on bad status")
	}
}

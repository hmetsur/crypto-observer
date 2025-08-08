package coingecko

import (
	"testing"
)

func TestGetPriceUSD(t *testing.T) {
	price, err := GetPriceUSD("btc")
	if err != nil {
		t.Fatalf("GetPriceUSD error: %v", err)
	}
	if price <= 0 {
		t.Errorf("GetPriceUSD: ожидал цену > 0, получил %.2f", price)
	}
}

// internal/coingecko/client.go
package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client получает цену в USD и возвращает её в центах
type Client struct {
	base string
	http *http.Client
}

func New(base string, timeout time.Duration) *Client {
	return &Client{
		base: base,
		http: &http.Client{Timeout: timeout},
	}
}

func (c *Client) GetPriceCents(ctx context.Context, symbol string) (int64, error) {
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.base, symbol)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var m map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return 0, err
	}
	usd := m[symbol]["usd"]
	return int64(usd * 100), nil
}

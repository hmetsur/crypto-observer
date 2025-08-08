package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var symbolToID = map[string]string{
	"btc":  "bitcoin",
	"eth":  "ethereum",
	"doge": "dogecoin",
}

func GetPriceUSD(symbol string) (float64, error) {
	id, ok := symbolToID[strings.ToLower(symbol)]
	if !ok {
		return 0, fmt.Errorf("неизвестный тикер: %s", symbol)
	}
	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd",
		id,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("статус: %d", resp.StatusCode)
	}
	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("ошибка декодирования: %w", err)
	}
	price, ok := data[id]["usd"]
	if !ok {
		return 0, fmt.Errorf("цена не найдена")
	}
	return price, nil
}

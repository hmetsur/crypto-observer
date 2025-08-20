package service

import "strings"

// Базовая карта популярных монет
var cgID = map[string]string{
	"btc":  "bitcoin",
	"eth":  "ethereum",
	"bnb":  "binancecoin",
	"sol":  "solana",
	"xrp":  "ripple",
	"ada":  "cardano",
	"doge": "dogecoin",
	"ton":  "the-open-network",
	"dot":  "polkadot",
	"trx":  "tron",
}

func toCoingeckoID(sym string) string {
	s := strings.ToLower(strings.TrimSpace(sym))
	if id, ok := cgID[s]; ok {
		return id
	}
	// Фоллбэк: пробуем как есть (вдруг уже совпадает с ID CoinGecko)
	return s
}

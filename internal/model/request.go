package model

type AddCurrencyRequest struct {
	Symbol string `json:"symbol"`
	Period int    `json:"period"`
}

type RemoveCurrencyRequest struct {
	Symbol string `json:"symbol"`
}

type PriceResponse struct {
	Coin      string  `json:"coin"`
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

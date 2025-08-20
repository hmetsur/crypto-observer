package model

type Price struct {
	Symbol string
	TS     int64
	Price  int64
}

type PriceDTO struct {
	Coin      string `json:"coin"`
	Timestamp int64  `json:"timestamp"`
	Price     int64  `json:"price"`
}

type AddReq struct {
	Symbol string `json:"symbol"` // например: "btc"
	Period int    `json:"period"` // период опроса в секундах
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func StatusOK() map[string]string {
	return map[string]string{"status": "ok"}
}

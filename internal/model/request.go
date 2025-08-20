package model

type AddCurrencyRequest struct {
	Symbol string `json:"symbol"`
	Period int    `json:"period"`
}

type RemoveReq struct {
	Symbol string `json:"symbol"`
}

type PriceResponse struct {
	Symbol string `json:"coin"`
	TS     int64  `json:"timestamp"`
	Price  int64  `json:"price"`
}

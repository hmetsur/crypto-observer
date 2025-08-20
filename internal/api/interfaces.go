package api

import "crypto-observer/internal/model"

type CurrencyService interface {
	AddCurrency(symbol string, periodSec int) error
	RemoveCurrency(symbol string) error
	GetPrice(symbol string, ts int64) (*model.Price, error)
}

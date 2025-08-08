package service

import (
	"errors"
	"fmt"
	"strconv"

	"crypto-observer/internal/db"
	"crypto-observer/internal/logger"
	"crypto-observer/internal/model"
)

type Collector interface {
	Start(symbol string, period int)
	Stop(symbol string)
}

type CurrencyService interface {
	AddCurrency(symbol string, period int) error
	RemoveCurrency(symbol string) error
	GetPrice(symbol, ts string) (*model.PriceResponse, error)
}

type Service struct {
	storage   db.StorageInterface
	collector Collector
}

var _ CurrencyService = (*Service)(nil)

func NewService(s db.StorageInterface, c Collector) *Service {
	return &Service{storage: s, collector: c}
}

func (s *Service) AddCurrency(symbol string, period int) error {
	logger.Log.WithFields(logger.Fields{"symbol": symbol, "period": period}).Info("Service: AddCurrency")
	s.collector.Start(symbol, period)
	return nil
}

func (s *Service) RemoveCurrency(symbol string) error {
	logger.Log.WithField("symbol", symbol).Info("Service: RemoveCurrency")
	s.collector.Stop(symbol)
	return nil
}

var (
	ErrPriceNotFound    = errors.New("price not found")
	ErrInvalidTimestamp = errors.New("invalid timestamp")
)

func (s *Service) GetPrice(symbol, ts string) (*model.PriceResponse, error) {
	logger.Log.WithFields(logger.Fields{"symbol": symbol, "timestamp": ts}).Info("Service: GetPrice")

	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		logger.Log.WithError(err).Warn("Service: bad timestamp")
		// заворачиваем исходную ошибку в свой сентинел
		return nil, fmt.Errorf("%w: %v", ErrInvalidTimestamp, err)
	}

	p, found, err := s.storage.GetClosestPrice(symbol, timestamp)
	if err != nil {
		logger.Log.WithError(err).Error("Service: storage error")
		return nil, err
	}
	if !found {
		logger.Log.WithField("symbol", symbol).Warn("Service: price not found")
		return nil, ErrPriceNotFound
	}

	return &model.PriceResponse{Coin: symbol, Timestamp: p.Timestamp, Price: p.Price}, nil
}

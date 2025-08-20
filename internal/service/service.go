// internal/service/service.go
package service

import (
	"context"
	"time"

	"crypto-observer/internal/coingecko"
	"crypto-observer/internal/model"
	"crypto-observer/pkg/logger"
)

type Storage interface {
	SavePrice(ctx context.Context, p model.Price) error
	GetClosestPrice(ctx context.Context, symbol string, ts int64) (*model.Price, error)
}

type Service struct {
	st         Storage
	collectors map[string]*collector
	defaultPer int
	priceCli   *coingecko.Client
}

func NewService(st Storage, defaultPeriod int, cgBaseURL string, timeout time.Duration) *Service {
	return &Service{
		st:         st,
		collectors: make(map[string]*collector),
		defaultPer: defaultPeriod,
		priceCli:   coingecko.New(cgBaseURL, timeout),
	}
}

func (s *Service) AddCurrency(symbol string, periodSec int) error {
	if periodSec <= 0 {
		periodSec = s.defaultPer
	}
	if c, ok := s.collectors[symbol]; ok && c.Running() {
		return nil
	}
	c := newCollector(symbol, time.Duration(periodSec)*time.Second, s.st, s.priceCli)
	s.collectors[symbol] = c
	c.Start()
	logger.L().WithField("symbol", symbol).Info("Service: AddCurrency")
	return nil
}

func (s *Service) RemoveCurrency(symbol string) error {
	if c, ok := s.collectors[symbol]; ok {
		c.Stop()
	}
	logger.L().WithField("symbol", symbol).Info("Service: RemoveCurrency")
	return nil
}

func (s *Service) GetPrice(symbol string, ts int64) (*model.Price, error) {
	logger.L().WithFields(logger.Fields{
		"symbol": symbol,
		"ts":     ts,
	}).Info("Service: GetPrice")

	// если ts == 0 — используем текущий момент
	if ts == 0 {
		ts = time.Now().Unix()
	}
	return s.st.GetClosestPrice(context.Background(), symbol, ts)
}

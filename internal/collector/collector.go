package collector

import (
	"context"
	"sync"
	"time"

	"crypto-observer/internal/coingecko"
	"crypto-observer/internal/db"
	"crypto-observer/internal/logger"
)

// Тип функции, которую можно подставить вместо реального CoinGecko
type GetPriceFunc func(symbol string) (float64, error)

type Collector interface {
	Start(symbol string, period int)
	Stop(symbol string)
}

type CollectorImpl struct {
	mu      sync.Mutex
	jobs    map[string]context.CancelFunc
	storage db.StorageInterface

	getPrice GetPriceFunc
}

func NewCollector(storage db.StorageInterface, f GetPriceFunc) *CollectorImpl {
	if f == nil {
		f = coingecko.GetPriceUSD
	}
	return &CollectorImpl{
		jobs:     make(map[string]context.CancelFunc),
		storage:  storage,
		getPrice: f,
	}
}

func (c *CollectorImpl) Start(symbol string, period int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.jobs[symbol]; ok {
		logger.Log.WithField("symbol", symbol).Warn("Collector: already running")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.jobs[symbol] = cancel

	logger.Log.WithFields(logger.Fields{"symbol": symbol, "period": period}).Info("Collector: start")

	go func() {
		ticker := time.NewTicker(time.Duration(period) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				price, err := c.getPrice(symbol)
				if err != nil {
					logger.Log.WithFields(logger.Fields{"symbol": symbol}).WithError(err).Error("Collector: fetch failed")
					continue
				}
				ts := time.Now().Unix()
				if err := c.storage.SavePrice(symbol, ts, price); err != nil {
					logger.Log.WithFields(logger.Fields{"symbol": symbol, "ts": ts, "price": price}).WithError(err).Error("Collector: save failed")
					continue
				}
				logger.Log.WithFields(logger.Fields{"symbol": symbol, "ts": ts, "price": price}).Info("Collector: saved")
			case <-ctx.Done():
				logger.Log.WithField("symbol", symbol).Info("Collector: stop")
				return
			}
		}
	}()
}

func (c *CollectorImpl) Running(symbol string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.jobs[symbol]
	return ok
}

func (c *CollectorImpl) Stop(symbol string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cancel, ok := c.jobs[symbol]; ok {
		cancel()
		delete(c.jobs, symbol)
		logger.Log.WithField("symbol", symbol).Info("Collector: stop request")
	} else {
		logger.Log.WithField("symbol", symbol).Warn("Collector: stop ignored (not running)")
	}
}

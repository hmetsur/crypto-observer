package service

import (
	"context"
	"sync/atomic"
	"time"

	"crypto-observer/internal/model"
	"crypto-observer/pkg/logger"
)

type priceClient interface {
	GetPriceCents(ctx context.Context, symbol string) (int64, error)
}

type storageIface interface {
	SavePrice(ctx context.Context, p model.Price) error
}

type collector struct {
	symbol string
	every  time.Duration
	st     storageIface
	pc     priceClient

	stopCh chan struct{}
	run    atomic.Bool // потокобезопасный флаг
}

func newCollector(symbol string, every time.Duration, st storageIface, pc priceClient) *collector {
	return &collector{
		symbol: symbol,
		every:  every,
		st:     st,
		pc:     pc,
		stopCh: make(chan struct{}, 1),
	}
}

func (c *collector) Start() {
	if c.run.Swap(true) { // если уже true — был запущен
		return
	}
	log := logger.L().WithField("symbol", c.symbol)
	go func() {
		t := time.NewTicker(c.every)
		defer t.Stop()
		defer log.Info("Collector: stop")
		for {
			select {
			case <-t.C:
				price, err := c.pc.GetPriceCents(context.Background(), toCoingeckoID(c.symbol))
				if err != nil {
					log.WithError(err).Error("Collector: fetch failed")
					continue
				}
				p := model.Price{Symbol: c.symbol, TS: time.Now().Unix(), Price: price}
				if err := c.st.SavePrice(context.Background(), p); err != nil {
					log.WithError(err).Error("Collector: save failed")
				}
			case <-c.stopCh:
				c.run.Store(false)
				return
			}
		}
	}()
	log.Info("Collector: start")
}

func (c *collector) Stop() {
	if !c.run.Load() {
		return
	}
	select {
	case c.stopCh <- struct{}{}:
	default:
	}
}

func (c *collector) Running() bool { return c.run.Load() }

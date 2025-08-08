package collector

import (
	"crypto-observer/internal/model"
	"testing"
	"time"
)

// Сигнализируем через канал, что SavePrice был вызван
type fakeStorage struct {
	called chan struct{}
}

func (f *fakeStorage) SavePrice(symbol string, ts int64, price float64) error {

	select {
	case f.called <- struct{}{}:
	default:
	}
	return nil
}

func (f *fakeStorage) GetClosestPrice(symbol string, ts int64) (*model.Price, bool, error) {
	return &model.Price{Symbol: symbol, Timestamp: ts, Price: 100.0}, true, nil
}

func TestCollector_StartAndStop(t *testing.T) {

	mockPrice := func(symbol string) (float64, error) {
		return 123.45, nil
	}

	fs := &fakeStorage{called: make(chan struct{}, 1)}

	coll := NewCollector(fs, mockPrice)

	coll.Start("btc", 1)
	defer coll.Stop("btc")

	select {
	case <-fs.called:

	case <-time.After(1500 * time.Millisecond):
		t.Fatalf("Ожидалось, что SavePrice будет вызван в течение ~1.5s, но этого не произошло")
	}

	coll.Stop("btc")
}

package collector

import (
	"testing"
	"time"

	"crypto-observer/internal/model"
)

// локальный фейк-хранилище (соответствует db.StorageInterface)
type fakeStorageRun struct {
	called chan struct{}
}

func (f *fakeStorageRun) SavePrice(symbol string, ts int64, price float64) error {
	// сигнализируем, что сохранение вызвалось
	select {
	case f.called <- struct{}{}:
	default:
	}
	return nil
}
func (f *fakeStorageRun) GetClosestPrice(symbol string, ts int64) (*model.Price, bool, error) {
	return &model.Price{Symbol: symbol, Timestamp: ts, Price: 1}, true, nil
}

func TestCollector_Running(t *testing.T) {
	fs := &fakeStorageRun{called: make(chan struct{}, 1)}
	mock := func(string) (float64, error) { return 1, nil }

	c := NewCollector(fs, mock)

	if c.Running("btc") {
		t.Fatalf("не должно быть запущено до Start")
	}
	c.Start("btc", 1)
	if !c.Running("btc") {
		t.Fatalf("ожидалось: запущено после Start")
	}

	select {
	case <-fs.called:
	case <-time.After(1500 * time.Millisecond):
		t.Fatalf("SavePrice не вызвался в разумное время")
	}

	c.Stop("btc")
	if c.Running("btc") {
		t.Fatalf("ожидалось: остановлено после Stop")
	}
	// повторный Stop — просто проверка, что не паникует
	c.Stop("btc")
}

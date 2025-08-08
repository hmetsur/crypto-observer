package service

import (
	"testing"

	"crypto-observer/internal/collector"
	"crypto-observer/internal/model"
)

// ---- фейки хранилища ----

type fakeStorage struct {
	prices []model.Price
}

func (f *fakeStorage) SavePrice(symbol string, ts int64, price float64) error {
	f.prices = append(f.prices, model.Price{Symbol: symbol, Timestamp: ts, Price: price})
	return nil
}

func (f *fakeStorage) GetClosestPrice(symbol string, ts int64) (*model.Price, bool, error) {
	for _, p := range f.prices {
		if p.Symbol == symbol && p.Timestamp <= ts {
			return &p, true, nil
		}
	}
	return nil, false, nil
}

// ---- тесты ----

func TestService_AddCurrency(t *testing.T) {
	fs := &fakeStorage{}
	mockPrice := func(symbol string) (float64, error) { return 999.99, nil }

	coll := collector.NewCollector(fs, mockPrice) // новая сигнатура (storage, getPriceFunc)
	s := NewService(fs, coll)

	if err := s.AddCurrency("btc", 1); err != nil {
		t.Fatalf("AddCurrency error: %v", err)
	}

	if !coll.Running("btc") {
		t.Errorf("ожидалось, что для btc будет запущен сбор цен")
	}
}

func TestService_RemoveCurrency(t *testing.T) {
	fs := &fakeStorage{}
	mockPrice := func(symbol string) (float64, error) { return 123.45, nil }

	coll := collector.NewCollector(fs, mockPrice)
	s := NewService(fs, coll)

	if err := s.AddCurrency("eth", 1); err != nil {
		t.Fatalf("prepare AddCurrency error: %v", err)
	}
	if !coll.Running("eth") {
		t.Fatalf("подготовка: ожидалось, что eth запущен")
	}

	if err := s.RemoveCurrency("eth"); err != nil {
		t.Fatalf("RemoveCurrency error: %v", err)
	}
	if coll.Running("eth") {
		t.Errorf("ожидалось, что для eth сбор будет остановлен")
	}
}

func TestService_GetPrice(t *testing.T) {
	// заранее положим запись, как будто коллектор уже писал
	fs := &fakeStorage{
		prices: []model.Price{
			{Symbol: "btc", Timestamp: 111, Price: 50000},
		},
	}
	mockPrice := func(symbol string) (float64, error) { return 50000, nil }

	coll := collector.NewCollector(fs, mockPrice)
	s := NewService(fs, coll)

	resp, err := s.GetPrice("btc", "111")
	if err != nil {
		t.Fatalf("GetPrice error: %v", err)
	}
	if resp == nil {
		t.Fatalf("ожидался непустой ответ")
	}
	if resp.Price != 50000 {
		t.Errorf("ожидалось 50000, получили %v", resp.Price)
	}
	if resp.Timestamp != 111 {
		t.Errorf("ожидался ts=111, получили %v", resp.Timestamp)
	}
}

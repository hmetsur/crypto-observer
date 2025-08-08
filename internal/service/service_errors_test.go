package service

import (
	"testing"

	"crypto-observer/internal/collector"
	"crypto-observer/internal/model"
)

type fakeStorageErr struct{}

func (f *fakeStorageErr) SavePrice(string, int64, float64) error { return nil }
func (f *fakeStorageErr) GetClosestPrice(string, int64) (*model.Price, bool, error) {
	return nil, false, nil
}

func TestService_GetPrice_BadTimestamp(t *testing.T) {
	fs := &fakeStorageErr{}
	c := collector.NewCollector(fs, func(string) (float64, error) { return 1, nil })
	s := NewService(fs, c)

	if _, err := s.GetPrice("btc", "abc"); err == nil {
		t.Fatalf("ожидалась ошибка парсинга timestamp")
	}
}

func TestService_GetPrice_NotFound(t *testing.T) {
	fs := &fakeStorageErr{}
	c := collector.NewCollector(fs, nil)
	s := NewService(fs, c)

	if _, err := s.GetPrice("btc", "123"); err == nil {
		t.Fatalf("ожидалась ошибка (not found)")
	}
}

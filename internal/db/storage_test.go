package db

import (
	"os"
	"testing"
	"time"
)

func TestStorage_SaveAndGetPrice(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN не задан, тест пропущен")
	}
	storage, err := NewStorage(dsn)
	if err != nil {
		t.Fatal("Ошибка подключения к БД:", err)
	}
	symbol := "btc"
	ts := time.Now().Unix()
	price := 12345.6

	err = storage.SavePrice(symbol, ts, price)
	if err != nil {
		t.Fatalf("Ошибка SavePrice: %v", err)
	}

	p, found, err := storage.GetClosestPrice(symbol, ts)
	if err != nil || !found {
		t.Fatalf("Ошибка GetClosestPrice: %v", err)
	}
	if p.Price != price {
		t.Errorf("Ожидал %.2f, получил %.2f", price, p.Price)
	}
}

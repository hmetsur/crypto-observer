//go:build !short

package db

import (
	"database/sql"
	"os"
	"testing"

	"crypto-observer/internal/model"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// Создаём схему, если её ещё нет
func setupDB(t *testing.T, db *sql.DB) {
	t.Helper()
	const schema = `
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    timestamp BIGINT NOT NULL,
    price DOUBLE PRECISION NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_prices_symbol_ts ON prices(symbol, timestamp);
`
	_, err := db.Exec(schema)
	require.NoError(t, err, "не удалось создать таблицу prices")
}

func TestStorage_SaveAndGetPrice_Full(t *testing.T) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		t.Skip("DB_DSN не задан, пропуск интеграционного теста")
	}

	// Подключаемся напрямую, чтобы выполнить setup
	raw, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer raw.Close()

	setupDB(t, raw)

	// Используем наш Storage поверх уже инициализированного *sql.DB
	st := &Storage{db: raw}

	// Сохраняем цену.
	symbol := "btc"
	ts := int64(111)
	val := 123.45

	err = st.SavePrice(symbol, ts, val)
	require.NoError(t, err, "SavePrice")

	// Читаем ближайшую цену
	got, found, err := st.GetClosestPrice(symbol, ts)
	require.NoError(t, err, "GetClosestPrice")
	require.True(t, found, "ожидалось found=true")
	require.NotNil(t, got, "результат должен быть не nil")

	// Проверяем поля структуры
	require.Equal(t, &model.Price{
		Symbol:    symbol,
		Timestamp: ts,
		Price:     val,
	}, got)
}

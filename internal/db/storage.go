package db

import (
	"database/sql"

	"crypto-observer/internal/logger"
	"crypto-observer/internal/model"

	_ "github.com/lib/pq"
)

type StorageInterface interface {
	SavePrice(symbol string, ts int64, price float64) error
	GetClosestPrice(symbol string, ts int64) (*model.Price, bool, error)
}

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := EnsureSchema(db); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func EnsureSchema(db *sql.DB) error {
	const schema = `
CREATE TABLE IF NOT EXISTS prices (
	id SERIAL PRIMARY KEY,
	symbol TEXT NOT NULL,
	timestamp BIGINT NOT NULL,
	price DOUBLE PRECISION NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_prices_symbol_ts ON prices(symbol, timestamp);
`
	if _, err := db.Exec(schema); err != nil {
		logger.Log.WithError(err).Error("DB: EnsureSchema failed")
		return err
	}
	logger.Log.Info("DB: schema ensured")
	return nil
}

func (s *Storage) SavePrice(symbol string, ts int64, price float64) error {
	logger.Log.WithFields(logger.Fields{"symbol": symbol, "ts": ts, "price": price}).Debug("DB: SavePrice")
	_, err := s.db.Exec(`INSERT INTO prices (symbol, timestamp, price) VALUES ($1, $2, $3)`, symbol, ts, price)
	if err != nil {
		logger.Log.WithError(err).Error("DB: SavePrice failed")
	}
	return err
}

func (s *Storage) GetClosestPrice(symbol string, ts int64) (*model.Price, bool, error) {
	logger.Log.WithFields(logger.Fields{"symbol": symbol, "ts": ts}).Debug("DB: GetClosestPrice")
	row := s.db.QueryRow(
		`SELECT symbol, timestamp, price
		 FROM prices
		 WHERE symbol = $1 AND timestamp <= $2
		 ORDER BY timestamp DESC
		 LIMIT 1`, symbol, ts,
	)
	var p model.Price
	if err := row.Scan(&p.Symbol, &p.Timestamp, &p.Price); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		logger.Log.WithError(err).Error("DB: GetClosestPrice failed")
		return nil, false, err
	}
	return &p, true, nil
}

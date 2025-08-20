// internal/db/storage.go
package db

import (
	"context"
	"errors"

	"crypto-observer/internal/model"
	"crypto-observer/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type poolIface interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}

type Storage struct {
	pool poolIface
}

func NewStorage(dsn string) (*Storage, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	st := &Storage{pool: pool}
	if err := st.EnsureSchema(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	return st, nil
}

func newWithPool(p poolIface) *Storage { return &Storage{pool: p} }

func (s *Storage) EnsureSchema(ctx context.Context) error {
	const q = `
CREATE TABLE IF NOT EXISTS prices (
    id           BIGSERIAL PRIMARY KEY,
    symbol       VARCHAR(32) NOT NULL,
    ts           BIGINT      NOT NULL,
    price_cents  BIGINT      NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_prices_symbol_ts ON prices(symbol, ts DESC);
`
	_, err := s.pool.Exec(ctx, q)
	if err != nil {
		logger.L().WithError(err).Error("DB: EnsureSchema failed")
	}
	return err
}

func (s *Storage) SavePrice(ctx context.Context, p model.Price) error {
	const q = `INSERT INTO prices (symbol, ts, price_cents) VALUES ($1, $2, $3)`
	_, err := s.pool.Exec(ctx, q, p.Symbol, p.TS, p.Price)
	if err != nil {
		logger.L().WithError(err).Error("DB: SavePrice failed")
	}
	return err
}

func (s *Storage) GetClosestPrice(ctx context.Context, symbol string, ts int64) (*model.Price, error) {
	const q = `
SELECT symbol, ts, price_cents
FROM prices
WHERE symbol = $1 AND ts <= $2
ORDER BY ts DESC
LIMIT 1`
	row := s.pool.QueryRow(ctx, q, symbol, ts)

	var out model.Price
	if err := row.Scan(&out.Symbol, &out.TS, &out.Price); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.L().WithError(err).Error("DB: GetClosestPrice failed")
		return nil, err
	}
	return &out, nil
}

func (s *Storage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

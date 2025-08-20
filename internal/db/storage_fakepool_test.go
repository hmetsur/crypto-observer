package db

import (
	"context"
	"errors"
	"testing"

	"crypto-observer/internal/model"
	"crypto-observer/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

// Инициализируем глобальный логгер, иначе logger.L() == nil и паникуем на WithError(...)
func init() { logger.Init() }

/*
 Тестируем Storage через внутренний конструктор newWithPool(p poolIface),
 подсовывая фейковый pool. Быстрые unit-тесты без реальной БД.
*/

type fakePool struct {
	execErr error
	row     pgx.Row
}

func (p *fakePool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, p.execErr
}
func (p *fakePool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, errors.New("not used")
}
func (p *fakePool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if p.row == nil {
		return fakeRow{scan: func(dest ...any) error { return pgx.ErrNoRows }}
	}
	return p.row
}
func (p *fakePool) Close() {}

type fakeRow struct {
	scan func(dest ...any) error
}

func (r fakeRow) Scan(dest ...any) error { return r.scan(dest...) }

func TestStorage_EnsureSchema_OK(t *testing.T) {
	fp := &fakePool{execErr: nil}
	st := newWithPool(fp)

	err := st.EnsureSchema(context.Background())
	require.NoError(t, err)
}

func TestStorage_SavePrice_OK(t *testing.T) {
	fp := &fakePool{execErr: nil}
	st := newWithPool(fp)

	err := st.SavePrice(context.Background(), model.Price{
		Symbol: "btc",
		TS:     111,
		Price:  12345,
	})
	require.NoError(t, err)
}

func TestStorage_GetClosestPrice_Found(t *testing.T) {
	row := fakeRow{
		scan: func(dest ...any) error {
			*(dest[0].(*string)) = "btc"
			*(dest[1].(*int64)) = 222
			*(dest[2].(*int64)) = 23456
			return nil
		},
	}
	fp := &fakePool{row: row}
	st := newWithPool(fp)

	got, err := st.GetClosestPrice(context.Background(), "btc", 999)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "btc", got.Symbol)
	require.Equal(t, int64(222), got.TS)
	require.Equal(t, int64(23456), got.Price)
}

func TestStorage_GetClosestPrice_NotFound(t *testing.T) {
	row := fakeRow{scan: func(dest ...any) error { return pgx.ErrNoRows }}
	fp := &fakePool{row: row}
	st := newWithPool(fp)

	got, err := st.GetClosestPrice(context.Background(), "btc", 999)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestStorage_GetClosestPrice_DBError(t *testing.T) {
	// вернём произвольную ошибку из Scan
	row := fakeRow{scan: func(dest ...any) error { return errors.New("db boom") }}
	fp := &fakePool{row: row}
	st := newWithPool(fp)

	got, err := st.GetClosestPrice(context.Background(), "btc", 100)
	require.Error(t, err)
	require.Nil(t, got)
}

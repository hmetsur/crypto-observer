package db

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"crypto-observer/internal/model"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// helper: создаём Storage с подменённым *sql.DB
func newMockStorage(t *testing.T) (*Storage, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, m, err := sqlmock.New() // sqlDB: *sql.DB
	require.NoError(t, err)

	st := &Storage{db: sqlDB}

	cleanup := func() {
		_ = st.db.Close() // без ExpectClose(), чтобы не ломать порядок ожиданий
	}

	return st, m, cleanup
}

func TestSavePrice_OK(t *testing.T) {
	st, m, done := newMockStorage(t)
	defer done()

	insert := regexp.QuoteMeta(`INSERT INTO prices (symbol, timestamp, price) VALUES ($1, $2, $3)`)

	m.ExpectExec(insert).
		WithArgs("btc", int64(111), 123.45).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := st.SavePrice("btc", 111, 123.45)
	require.NoError(t, err)

	require.NoError(t, m.ExpectationsWereMet())
}

func TestSavePrice_DBError(t *testing.T) {
	st, m, done := newMockStorage(t)
	defer done()

	insert := regexp.QuoteMeta(`INSERT INTO prices (symbol, timestamp, price) VALUES ($1, $2, $3)`)

	m.ExpectExec(insert).
		WithArgs("btc", int64(111), 123.45).
		WillReturnError(errors.New("db failed"))

	err := st.SavePrice("btc", 111, 123.45)
	require.Error(t, err)

	require.NoError(t, m.ExpectationsWereMet())
}

func TestGetClosestPrice_Found(t *testing.T) {
	st, m, done := newMockStorage(t)
	defer done()

	const selectPricesRe = `(?i)SELECT\s+symbol,\s+timestamp,\s+price\s+` +
		`FROM\s+prices\s+` +
		`WHERE\s+symbol\s*=\s*\$1\s+AND\s+timestamp\s*<=\s*\$2\s+` +
		`ORDER\s+BY\s+timestamp\s+DESC\s+LIMIT\s+1`

	rows := sqlmock.NewRows([]string{"symbol", "timestamp", "price"}).
		AddRow("btc", int64(111), 123.45)

	m.ExpectQuery(selectPricesRe).
		WithArgs("btc", int64(200)).
		WillReturnRows(rows)

	res, found, err := st.GetClosestPrice("btc", 200)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, &model.Price{Symbol: "btc", Timestamp: 111, Price: 123.45}, res)

	require.NoError(t, m.ExpectationsWereMet())
}

func TestGetClosestPrice_NoRows(t *testing.T) {
	st, m, done := newMockStorage(t)
	defer done()

	const selectPricesRe = `(?i)SELECT\s+symbol,\s+timestamp,\s+price\s+` +
		`FROM\s+prices\s+` +
		`WHERE\s+symbol\s*=\s*\$1\s+AND\s+timestamp\s*<=\s*\$2\s+` +
		`ORDER\s+BY\s+timestamp\s+DESC\s+LIMIT\s+1`

	m.ExpectQuery(selectPricesRe).
		WithArgs("eth", int64(999)).
		WillReturnError(sql.ErrNoRows)

	res, found, err := st.GetClosestPrice("eth", 999)
	require.NoError(t, err)
	require.False(t, found)
	require.Nil(t, res)

	require.NoError(t, m.ExpectationsWereMet())
}

func TestGetClosestPrice_DBError(t *testing.T) {
	st, m, done := newMockStorage(t)
	defer done()

	const selectPricesRe = `(?i)SELECT\s+symbol,\s+timestamp,\s+price\s+` +
		`FROM\s+prices\s+` +
		`WHERE\s+symbol\s*=\s*\$1\s+AND\s+timestamp\s*<=\s*\$2\s+` +
		`ORDER\s+BY\s+timestamp\s+DESC\s+LIMIT\s+1`

	m.ExpectQuery(selectPricesRe).
		WithArgs("btc", int64(100)).
		WillReturnError(errors.New("db boom"))

	res, found, err := st.GetClosestPrice("btc", 100)
	require.Error(t, err)
	require.False(t, found)
	require.Nil(t, res)

	require.NoError(t, m.ExpectationsWereMet())
}

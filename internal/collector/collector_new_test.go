package collector

import (
	"testing"

	"crypto-observer/internal/model"

	"github.com/stretchr/testify/require"
)

// мок StorageInterface
type dummyStorage struct{}

func (d *dummyStorage) SavePrice(symbol string, ts int64, price float64) error { return nil }

func (d *dummyStorage) GetClosestPrice(symbol string, ts int64) (*model.Price, bool, error) {
	return &model.Price{Symbol: symbol, Timestamp: ts, Price: 100}, true, nil
}

func TestNewCollector(t *testing.T) {
	st := &dummyStorage{}
	c := NewCollector(st, func(symbol string) (float64, error) { return 100.0, nil })
	require.NotNil(t, c)
}

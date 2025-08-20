package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"crypto-observer/internal/model"

	"github.com/stretchr/testify/require"
)

// ---- fakes ----

type fakeStorage struct {
	mu        sync.Mutex
	gotSym    string
	gotTS     int64
	retPrice  *model.Price
	retErr    error
	saveCalls int
}

func (f *fakeStorage) SavePrice(ctx context.Context, p model.Price) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.saveCalls++
	return nil
}

func (f *fakeStorage) GetClosestPrice(ctx context.Context, symbol string, ts int64) (*model.Price, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.gotSym = symbol
	f.gotTS = ts
	return f.retPrice, f.retErr
}

// ---- helpers ----

func newSvcWith(storage Storage) *Service {
	// baseURL/timeout не важны — в тестах не даём коллектору тикаать
	return NewService(storage /*defaultPeriod*/, 1, "http://localhost", 5*time.Second)
}

func sleepMS(ms int) { time.Sleep(time.Duration(ms) * time.Millisecond) }

// ---- tests ----

func TestService_GetPrice_ZeroTS_UsesNow(t *testing.T) {
	fs := &fakeStorage{retPrice: &model.Price{Symbol: "btc", TS: 1, Price: 2}}
	s := newSvcWith(fs)

	start := time.Now().Unix()
	got, err := s.GetPrice("btc", 0)
	require.NoError(t, err)
	require.Equal(t, fs.retPrice, got)

	// проверяем, что в сторадж ушёл ts "примерно сейчас"
	require.Equal(t, "btc", fs.gotSym)
	require.InDelta(t, start, fs.gotTS, 2, "ts should be near now (seconds)")
}

func TestService_GetPrice_PassesTS(t *testing.T) {
	fs := &fakeStorage{retPrice: &model.Price{Symbol: "eth", TS: 111, Price: 222}}
	s := newSvcWith(fs)

	got, err := s.GetPrice("eth", 12345)
	require.NoError(t, err)
	require.Equal(t, fs.retPrice, got)
	require.Equal(t, int64(12345), fs.gotTS)
	require.Equal(t, "eth", fs.gotSym)
}

func TestService_GetPrice_PropagatesError(t *testing.T) {
	fs := &fakeStorage{retErr: errors.New("db boom")}
	s := newSvcWith(fs)

	got, err := s.GetPrice("btc", 100)
	require.Error(t, err)
	require.Nil(t, got)
}

func TestService_AddCurrency_StartsCollector_AndRemoveStops(t *testing.T) {
	fs := &fakeStorage{}
	s := newSvcWith(fs)

	// длинный период — чтобы тики не успели сработать в тесте
	err := s.AddCurrency("btc", 3600)
	require.NoError(t, err)

	c, ok := s.collectors["btc"]
	require.True(t, ok, "collector must be created")
	require.True(t, c.Running(), "collector should be running after AddCurrency")

	// Теперь удаляем и убеждаемся, что остановился
	err = s.RemoveCurrency("btc")
	require.NoError(t, err)
	sleepMS(30)
	require.False(t, c.Running(), "collector should stop after RemoveCurrency")
}

func TestService_AddCurrency_DefaultPeriodUsed_WhenZero(t *testing.T) {
	fs := &fakeStorage{}
	s := newSvcWith(fs)
	// periodSec <= 0 => берётся defaultPer
	err := s.AddCurrency("eth", 0)
	require.NoError(t, err)
	c, ok := s.collectors["eth"]
	require.True(t, ok)
	require.True(t, c.Running())
	// уборка
	_ = s.RemoveCurrency("eth")
	sleepMS(20)
}

func TestService_AddCurrency_SecondCallDoesNotDuplicate(t *testing.T) {
	fs := &fakeStorage{}
	s := newSvcWith(fs)

	require.NoError(t, s.AddCurrency("btc", 3600))
	first := s.collectors["btc"]
	require.NotNil(t, first)

	// повторный вызов — коллектор уже запущен; должен остаться тем же
	require.NoError(t, s.AddCurrency("btc", 1))
	second := s.collectors["btc"]

	require.Same(t, first, second, "should not replace already running collector")
	_ = s.RemoveCurrency("btc")
	sleepMS(20)
}

package service

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"crypto-observer/internal/model"

	"crypto-observer/pkg/logger"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.Init()
}

type fakePriceClient struct {
	val   int64
	err   error
	calls int32
}

func (f *fakePriceClient) GetPriceCents(ctx context.Context, symbol string) (int64, error) {
	atomic.AddInt32(&f.calls, 1)
	return f.val, f.err
}

type memStorage struct {
	mu    sync.Mutex
	last  model.Price
	err   error
	count int32
}

func (m *memStorage) SavePrice(ctx context.Context, p model.Price) error {
	atomic.AddInt32(&m.count, 1)
	m.mu.Lock()
	m.last = p
	m.mu.Unlock()
	return m.err
}

// ----- helpers -----

func wait(ms int) { time.Sleep(time.Duration(ms) * time.Millisecond) }

// ----- tests -----

func TestCollector_StartAndStop_OK(t *testing.T) {
	st := &memStorage{}
	pc := &fakePriceClient{val: 12345}

	c := newCollector("btc", 60*time.Millisecond, st, pc)

	require.False(t, c.Running(), "should be not running before start")
	c.Start()
	require.True(t, c.Running(), "should be running after start")

	// ждём, чтобы тиковщик успел отработать хотя бы раз
	wait(150)

	c.Stop()
	// дать горутине завершиться
	wait(40)
	require.False(t, c.Running(), "should not be running after stop")

	require.GreaterOrEqual(t, atomic.LoadInt32(&st.count), int32(1), "expected at least one SavePrice call")
	st.mu.Lock()
	defer st.mu.Unlock()
	require.Equal(t, "btc", st.last.Symbol)
	require.Equal(t, pc.val, st.last.Price)
	require.NotZero(t, st.last.TS)
}

func TestCollector_DoubleStart_NoPanic(t *testing.T) {
	st := &memStorage{}
	pc := &fakePriceClient{val: 7}

	c := newCollector("eth", 50*time.Millisecond, st, pc)

	// второй старт должен быть проигнорирован и не паниковать
	c.Start()
	c.Start()
	wait(130)

	c.Stop()
	wait(30)

	require.GreaterOrEqual(t, atomic.LoadInt32(&st.count), int32(1))
}

func TestCollector_FetchError(t *testing.T) {
	st := &memStorage{}
	pc := &fakePriceClient{err: errors.New("boom")}

	c := newCollector("btc", 50*time.Millisecond, st, pc)
	c.Start()
	wait(140)
	c.Stop()
	wait(30)

	require.Equal(t, int32(0), atomic.LoadInt32(&st.count), "SavePrice shouldn't be called when fetch fails")
	require.GreaterOrEqual(t, atomic.LoadInt32(&pc.calls), int32(1), "GetPriceCents must be attempted")
}

func TestCollector_SaveError_StillTicks(t *testing.T) {
	st := &memStorage{err: errors.New("db-fail")}
	pc := &fakePriceClient{val: 999}

	c := newCollector("btc", 50*time.Millisecond, st, pc)
	c.Start()
	wait(140)
	c.Stop()
	wait(30)

	// Сохранение падает, но попытки должны происходить
	require.GreaterOrEqual(t, atomic.LoadInt32(&st.count), int32(1))
}

func TestCollector_RunningFlagTransitions(t *testing.T) {
	st := &memStorage{}
	pc := &fakePriceClient{val: 1}
	c := newCollector("btc", 80*time.Millisecond, st, pc)

	require.False(t, c.Running())
	c.Start()
	require.True(t, c.Running())
	c.Stop()
	wait(30)
	require.False(t, c.Running())
}

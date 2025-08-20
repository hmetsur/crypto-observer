package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-observer/internal/model"
	"crypto-observer/pkg/logger"
	"github.com/stretchr/testify/require"
)

type fakeService struct {
	addErr  error
	rmErr   error
	getResp *model.Price
	getErr  error
	gotAdd  struct {
		symbol string
		period int
	}
	gotRemove string
	gotGet    struct {
		symbol string
		ts     int64
	}
}

func (f *fakeService) AddCurrency(symbol string, period int) error {
	f.gotAdd = struct {
		symbol string
		period int
	}{symbol, period}
	return f.addErr
}
func (f *fakeService) RemoveCurrency(symbol string) error {
	f.gotRemove = symbol
	return f.rmErr
}
func (f *fakeService) GetPrice(symbol string, ts int64) (*model.Price, error) {
	f.gotGet = struct {
		symbol string
		ts     int64
	}{symbol, ts}
	return f.getResp, f.getErr
}

func init() { logger.Init() }

func TestHandler_AddCurrency(t *testing.T) {
	tests := []struct {
		name     string
		body     any
		svcErr   error
		wantCode int
	}{
		{"ok", map[string]any{"symbol": "btc", "period": 5}, nil, http.StatusOK},
		{"bad json", "not-json", nil, http.StatusBadRequest},
		{"missing fields", map[string]any{"symbol": ""}, nil, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := &fakeService{addErr: tc.svcErr}
			h := NewHandler(fs)

			var b []byte
			switch v := tc.body.(type) {
			case string:
				b = []byte(v)
			default:
				b, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/currency/add", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.AddCurrency(rr, req)
			require.Equal(t, tc.wantCode, rr.Code)
			if tc.wantCode == http.StatusOK {
				require.Equal(t, "btc", fs.gotAdd.symbol)
				require.Equal(t, 5, fs.gotAdd.period)
			}
		})
	}
}

func TestHandler_RemoveCurrency(t *testing.T) {
	tests := []struct {
		name     string
		body     any
		svcErr   error
		wantCode int
	}{
		{"ok", map[string]any{"symbol": "btc"}, nil, http.StatusOK},
		{"bad json", "oops", nil, http.StatusBadRequest},
		{"missing symbol", map[string]any{"symbol": ""}, nil, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := &fakeService{rmErr: tc.svcErr}
			h := NewHandler(fs)

			var b []byte
			switch v := tc.body.(type) {
			case string:
				b = []byte(v)
			default:
				b, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/currency/remove", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.RemoveCurrency(rr, req)
			require.Equal(t, tc.wantCode, rr.Code)
			if rr.Code == http.StatusOK {
				require.Equal(t, "btc", fs.gotRemove)
			}
		})
	}
}

func TestHandler_GetPrice(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		resp     *model.Price
		svcErr   error
		wantCode int
	}{
		{
			name:     "ok",
			query:    "/currency/price?symbol=btc&timestamp=111",
			resp:     &model.Price{Symbol: "btc", TS: 111, Price: 12345},
			wantCode: http.StatusOK,
		},
		{
			name:     "missing params",
			query:    "/currency/price?symbol=&timestamp=",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "bad ts",
			query:    "/currency/price?symbol=btc&timestamp=abc",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			query:    "/currency/price?symbol=btc&timestamp=111",
			resp:     nil,
			svcErr:   nil,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := &fakeService{getResp: tc.resp, getErr: tc.svcErr}
			h := NewHandler(fs)

			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			rr := httptest.NewRecorder()

			h.GetPrice(rr, req)
			require.Equal(t, tc.wantCode, rr.Code)
			if rr.Code == http.StatusOK {
				var out struct{ Price int64 }
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &out))
				require.Equal(t, int64(12345), out.Price)
			}
		})
	}
}

// простой маркер ошибки
type markerErr string

func (m markerErr) Error() string { return string(m) }
func assertError(s string) error  { return markerErr(s) }

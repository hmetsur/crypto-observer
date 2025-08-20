BEGIN;

CREATE TABLE IF NOT EXISTS prices (
    id      BIGSERIAL PRIMARY KEY,
    symbol  VARCHAR(32)   NOT NULL,
    ts      BIGINT        NOT NULL,
    price_cents   BIGINT        NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_prices_symbol_ts
    ON prices(symbol, ts DESC);

COMMIT;
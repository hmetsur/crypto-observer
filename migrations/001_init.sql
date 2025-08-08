CREATE TABLE prices (
                        id SERIAL PRIMARY KEY,
                        symbol TEXT NOT NULL,
                        timestamp BIGINT NOT NULL,
                        price DOUBLE PRECISION NOT NULL
);
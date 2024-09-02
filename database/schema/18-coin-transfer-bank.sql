-- +migrate Up
CREATE TABLE  coin_transfer
(
    id         SERIAL PRIMARY KEY,
    event      TEXT,
    address    TEXT,
    coins      COIN[] NOT NULL,
    height     BIGINT  NOT NULL
);

-- +migrate Down
DROP TABLE coin_transfer;

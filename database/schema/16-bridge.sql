-- +migrate Up
CREATE TABLE bridge_params
(
    one_row_id            BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    withdraw_denom        TEXT    NOT NULL,
    height                BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX bridge_params_height_index ON bridge_params (height);

CREATE TABLE hash
(
    index                   TEXT UNIQUE NOT NULL PRIMARY KEY REFERENCES operation(index)
);

-- +migrate Down
DROP TABLE hash;
DROP TABLE bridge_params;



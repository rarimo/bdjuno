-- +migrate Up
CREATE TYPE NETWORK_PARAMS as
(
    name    TEXT,
    contact TEXT,
    type    INT
);

CREATE TYPE TOKENMANAGER_PARAMS as
(
    networks NETWORK_PARAMS[]
);


CREATE TABLE tokenmanager_params
(
    one_row_id BOOLEAN             NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     TOKENMANAGER_PARAMS NOT NULL DEFAULT '{}'::TOKENMANAGER_PARAMS,
    height     BIGINT              NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX tokenmanager_params_height_index ON tokenmanager_params (height);


CREATE TYPE COLLECTION_DATA_INDEX as
(
    chain   TEXT,
    address TEXT
);

CREATE TYPE COLLECTION_METADATA as
(
    name         TEXT,
    symbol       TEXT,
    metadata_uri TEXT
);

CREATE TABLE collection
(
    index_key BYTEA                   NOT NULL PRIMARY KEY,
    index     TEXT                    NOT NULL,
    meta      COLLECTION_METADATA     NOT NULL,
    data      COLLECTION_DATA_INDEX[] NOT NULL DEFAULT '[]'::COLLECTION_DATA_INDEX[]
);


CREATE TABLE collection_data
(
    index_key  BYTEA                 NOT NULL PRIMARY KEY,
    index      COLLECTION_DATA_INDEX NOT NULL DEFAULT '{}'::COLLECTION_DATA_INDEX,
    collection TEXT REFERENCES collection (index),
    token_type INT                   NOT NULL,
    wrapped    BOOLEAN               NOT NULL,
    decimals   INT                   NOT NULL
);

CREATE TYPE ON_CHAIN_ITEM_INDEX as
(
    chain    TEXT,
    address  TEXT,
    token_id TEXT
);

CREATE TYPE ITEM_METADATA as
(
    image_uri  TEXT,
    image_hash TEXT,
    seed       TEXT,
    name       TEXT,
    symbol     TEXT,
    uri        TEXT
);

CREATE TABLE item
(
    index      TEXT                  NOT NULL PRIMARY KEY,
    collection TEXT                  NOT NULL REFERENCES collection (index),
    meta       ITEM_METADATA         NOT NULL DEFAULT '{}'::ITEM_METADATA,
    on_chain   ON_CHAIN_ITEM_INDEX[] NOT NULL DEFAULT '[]'::ON_CHAIN_ITEM_INDEX[]
);

CREATE TABLE on_chain_item
(
    index     ON_CHAIN_ITEM_INDEX NOT NULL,
    item      TEXT                NOT NULL REFERENCES item (index)
);

CREATE TABLE seed
(
    seed TEXT NOT NULL PRIMARY KEY,
    item TEXT NOT NULL REFERENCES item (index)
);

-- +migrate Down
DROP TABLE seed;
DROP TABLE on_chain_item;
DROP TABLE item;
DROP TYPE ITEM_METADATA;
DROP TYPE ON_CHAIN_ITEM_INDEX;
DROP TABLE collection_data;
DROP TABLE collection;
DROP TYPE COLLECTION_METADATA;
DROP TYPE COLLECTION_DATA_INDEX;
DROP TABLE tokenmanager_params;
DROP TYPE TOKENMANAGER_PARAMS;
DROP TYPE NETWORK_PARAMS;

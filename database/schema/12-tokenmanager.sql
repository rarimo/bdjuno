-- +migrate Up
CREATE TABLE network
(
    name   TEXT  NOT NULL PRIMARY KEY,
    type   INT   NOT NULL,
    params JSONB NOT NULL
);

CREATE TABLE collection
(
    index TEXT UNIQUE NOT NULL PRIMARY KEY,
    meta  JSONB       NOT NULL,
    data  JSONB       NOT NULL
);

CREATE TABLE collection_data
(
    index_key  TEXT    NOT NULL PRIMARY KEY,
    index      JSONB   NOT NULL,
    collection TEXT REFERENCES collection (index),
    token_type INT     NOT NULL,
    wrapped    BOOLEAN NOT NULL,
    decimals   INT     NOT NULL
);

CREATE TABLE item
(
    index      TEXT  NOT NULL PRIMARY KEY,
    collection TEXT  NOT NULL REFERENCES collection (index),
    meta       JSONB NOT NULL,
    on_chain   JSONB NOT NULL
);

CREATE TABLE on_chain_item
(
    index JSONB NOT NULL,
    item  TEXT  NOT NULL REFERENCES item (index)
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
DROP TABLE collection_data;
DROP TABLE collection;
DROP TABLE network;

-- +migrate Up
CREATE TABLE parties
(
    account  TEXT    NOT NULL PRIMARY KEY REFERENCES account (address),
    pub_key  TEXT    NOT NULL,
    address  TEXT    NOT NULL,
    verified BOOLEAN NOT NULL
);

CREATE TABLE rarimocore_params
(
    one_row_id         BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    key_ecdsa          TEXT    NOT NULL,
    threshold          BIGINT  NOT NULL,
    is_update_required BOOLEAN NOT NULL,
    last_signature     TEXT    NOT NULL,
    parties            TEXT[]  NOT NULL,
    height             BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX rarimocore_params_height_index ON rarimocore_params (height);

CREATE TABLE operation
(
    index          TEXT   NOT NULL PRIMARY KEY,
    operation_type INT    NOT NULL,
    status         INT    NOT NULL,
    creator        TEXT   NOT NULL REFERENCES account (address),
    timestamp      BIGINT NOT NULL
);

CREATE TABLE transfer
(
    operation_index TEXT  NOT NULL PRIMARY KEY REFERENCES operation (index),
    origin          TEXT  NOT NULL,
    tx              TEXT  NOT NULL,
    event_id        TEXT  NOT NULL,
    receiver        TEXT  NOT NULL,
    amount          TEXT  NOT NULL,
    bundle_data     TEXT,
    bundle_salt     TEXT,
    "from"          JSONB NOT NULL,
    "to"            JSONB NOT NULL,
    item_meta       JSONB -- Optional (if item currently does not exists)
);

CREATE TABLE change_parties
(
    operation_index TEXT   NOT NULL PRIMARY KEY REFERENCES operation (index),
    parties         TEXT[] NOT NULL,
    new_public_key  TEXT   NOT NULL,
    signature       TEXT   NOT NULL
);

CREATE TABLE confirmation
(
    root            TEXT   NOT NULL PRIMARY KEY,
    indexes         TEXT[] NOT NULL,
    signature_ecdsa TEXT   NOT NULL,
    creator         TEXT   NOT NULL REFERENCES account (address),
    height          BIGINT NOT NULL,
    tx              TEXT
);

CREATE TABLE vote
(
    operation TEXT   NOT NULL PRIMARY KEY REFERENCES operation (index),
    validator TEXT   NOT NULL REFERENCES account (address),
    vote      INT    NOT NULL,
    height    BIGINT NOT NULL,
    tx        TEXT
);

CREATE TABLE contract_upgrade
(
    operation_index             TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    target_contract             TEXT NOT NULL,
    chain                       TEXT NOT NULL,
    new_implementation_contract TEXT NOT NULL,
    hash                        TEXT NOT NULL,
    buffer_account              TEXT NOT NULL,
    nonce                       TEXT NOT NULL,
    type                        INT  NOT NULL
);

CREATE TABLE fee_token_management
(
    operation_index    TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    op_type            INT  NOT NULL,
    fee_token_contract TEXT NOT NULL,
    fee_token_amount   TEXT NOT NULL,
    chain              TEXT NOT NULL,
    receiver           TEXT NOT NULL,
    nonce              TEXT NOT NULL
);

CREATE TABLE identity_default_transfer
(
    operation_index            TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    contract                   TEXT NOT NULL,
    chain                      TEXT NOT NULL,
    gisthash                   TEXT NOT NULL,
    id                         TEXT NOT NULL,
    state_hash                 TEXT NOT NULL,
    state_created_at_timestamp TEXT NOT NULL,
    state_created_at_block     TEXT NOT NULL,
    state_replaced_by          TEXT NOT NULL,
    gistreplaced_by            TEXT NOT NULL,
    gistcreated_at_timestamp   TEXT NOT NULL,
    gistcreated_at_block       TEXT NOT NULL,
    replaced_state_hash        TEXT NOT NULL,
    replaced_gist_hash         TEXT NOT NULL
);

CREATE TABLE identity_gist_transfer
(
    operation_index            TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    contract                   TEXT NOT NULL,
    chain                      TEXT NOT NULL,
    gisthash                   TEXT NOT NULL,
    gistcreated_at_timestamp   TEXT NOT NULL,
    gistcreated_at_block       TEXT NOT NULL,
    replaced_gist_hash         TEXT NOT NULL
);

CREATE TABLE identity_state_transfer
(
    operation_index            TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    contract                   TEXT NOT NULL,
    chain                      TEXT NOT NULL,
    id                         TEXT NOT NULL,
    state_hash                 TEXT NOT NULL,
    state_created_at_timestamp TEXT NOT NULL,
    state_created_at_block     TEXT NOT NULL,
    replaced_state_hash        TEXT NOT NULL
);

-- +migrate Down
DROP TABLE identity_state_transfer;
DROP TABLE identity_gist_transfer;
DROP TABLE identity_default_transfer;
DROP TABLE fee_token_management;
DROP TABLE contract_upgrade;
DROP TABLE vote;
DROP TABLE confirmation;
DROP TABLE change_parties;
DROP TABLE transfer;
DROP TABLE operation;
DROP TABLE rarimocore_params;
DROP TABLE parties;

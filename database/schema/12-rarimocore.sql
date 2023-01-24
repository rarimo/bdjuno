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
    one_row_id                   BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    key_ecdsa                    TEXT    NOT NULL,
    threshold                    BIGINT  NOT NULL,
    is_update_required           BOOLEAN NOT NULL,
    last_signature               TEXT    NOT NULL,
    parties                      TEXT[]  NOT NULL DEFAULT '[]'::TEXT[],
    height                       BIGINT  NOT NULL,
    available_resign_block_delta BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX rarimocore_params_height_index ON rarimocore_params (height);

CREATE TABLE operation
(
    index          TEXT    NOT NULL PRIMARY KEY,
    operation_type INT     NOT NULL,
    signed         BOOLEAN NOT NULL,
    creator        TEXT    NOT NULL REFERENCES account (address),
    timestamp      BIGINT  NOT NULL
);

CREATE TABLE transfer
(
    operation_index TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    origin          TEXT NOT NULL,
    tx              TEXT NOT NULL,
    event_id        TEXT NOT NULL,
    from_chain      TEXT NOT NULL,
    to_chain        TEXT NOT NULL,
    receiver        TEXT NOT NULL,
    amount          TEXT NOT NULL,
    bundle_data     TEXT,
    bundle_salt     TEXT,
    token_index     TEXT NOT NULL -- TODO: ADD REFERENCES TO TOKEN TABLE
);

CREATE TABLE change_parties
(
    operation_index TEXT   NOT NULL PRIMARY KEY REFERENCES operation (index),
    parties         TEXT[] NOT NULL DEFAULT '[]'::TEXT[],
    new_public_key  TEXT   NOT NULL,
    signature       TEXT   NOT NULL
);

CREATE TABLE confirmation
(
    root            TEXT   NOT NULL PRIMARY KEY,
    indexes         TEXT[] NOT NULL DEFAULT '[]',
    signature_ecdsa TEXT   NOT NULL,
    creator         TEXT   NOT NULL REFERENCES account (address)
);

-- +migrate Down
DROP TABLE confirmation;
DROP TABLE change_parties;
DROP TABLE transfer;
DROP TABLE operation;
DROP TABLE rarimocore_params;
DROP TABLE parties;

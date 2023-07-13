-- +migrate Up
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
    operation_index             TEXT NOT NULL PRIMARY KEY REFERENCES operation (index),
    contract                    TEXT NOT NULL,
    chain                       TEXT NOT NULL,
    gisthash                    TEXT NOT NULL,
    id                          TEXT NOT NULL,
    state_hash                  TEXT NOT NULL,
    state_created_at_timestamp  TEXT NOT NULL,
    state_created_at_block      TEXT NOT NULL,
    state_replaced_by           TEXT NOT NULL,
    gistreplaced_by             TEXT NOT NULL,
    gistcreated_at_timestamp    TEXT NOT NULL,
    gistcreated_at_block        TEXT NOT NULL,
    replaced_state_hash         TEXT NOT NULL,
    replaced_gist_hash          TEXT NOT NULL
);

-- +migrate Down
DROP TABLE identity_default_transfer;
DROP TABLE fee_token_management;
DROP TABLE contract_upgrade;

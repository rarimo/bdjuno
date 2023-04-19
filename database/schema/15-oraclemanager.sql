-- +migrate Up
CREATE TABLE oraclemanager_params
(
    one_row_id            BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    min_oracle_stake      TEXT    NOT NULL,
    check_operation_delta INT     NOT NULL,
    max_violations_count  INT     NOT NULL,
    max_missed_count      INT     NOT NULL,
    slashed_freeze_blocks INT     NOT NULL,
    min_oracles_count     INT     NOT NULL,
    stake_denom           TEXT    NOT NULL,
    vote_quorum           TEXT    NOT NULL,
    vote_threshold        TEXT    NOT NULL,
    height                BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX oraclemanager_params_height_index ON oraclemanager_params (height);

CREATE TABLE oracle
(
    index                   TEXT UNIQUE NOT NULL PRIMARY KEY,
    chain                   TEXT        NOT NULL,
    account                 TEXT        NOT NULL REFERENCES account (address),
    status                  INT         NOT NULL,
    stake                   TEXT        NOT NULL,
    missed_count            INT         NOT NULL,
    violations_count        INT         NOT NULL,
    freeze_end_block        INT         NOT NULL,
    votes_count             INT         NOT NULL,
    create_operations_count INT         NOT NULL
);

-- +migrate Down
DROP TABLE oracle;
DROP TABLE oraclemanager_params;



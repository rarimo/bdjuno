-- +migrate Up
CREATE TABLE multisig_params
(
    one_row_id       BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    group_sequence    INT     NOT NULL,
    proposal_sequence INT     NOT NULL,
    prune_period      INT     NOT NULL,
    voting_period     TEXT    NOT NULL,
    height           BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX multisig_params_height_index ON multisig_params (height);

CREATE TABLE "group"
(
    account   TEXT UNIQUE NOT NULL PRIMARY KEY REFERENCES account (address),
    members   TEXT[]      NOT NULL,
    threshold INT         NOT NULL
);

CREATE TABLE multisig_proposal
(
    id                 INT UNIQUE NOT NULL PRIMARY KEY,
    proposer           TEXT       NOT NULL REFERENCES account (address),
    "group"            TEXT       NOT NULL REFERENCES "group" (account),
    submit_block       BIGINT     NOT NULL,
    voting_end_block   BIGINT     NOT NULL,
    status             TEXT       NOT NULL,
    final_tally_result JSONB      NOT NULL,
    messages           JSONB      NOT NULL
);

CREATE TABLE multisig_proposal_vote
(
    index        TEXT UNIQUE NOT NULL PRIMARY KEY,
    voter        TEXT        NOT NULL REFERENCES account (address),
    proposal_id  INT         NOT NULL REFERENCES multisig_proposal (id),
    option       INT         NOT NULL,
    submit_block BIGINT      NOT NULL
);

-- +migrate Down
DROP TABLE multisig_proposal_vote;
DROP TABLE multisig_proposal;
DROP TABLE "group";
DROP TABLE multisig_params;



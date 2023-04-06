-- +migrate Up
ALTER TABLE rarimocore_params
    ADD COLUMN IF NOT EXISTS stake_amount TEXT;
ALTER TABLE rarimocore_params
    ADD COLUMN IF NOT EXISTS stake_denom TEXT;
ALTER TABLE rarimocore_params
    ADD COLUMN IF NOT EXISTS max_violations_count INT;
ALTER TABLE rarimocore_params
    ADD COLUMN IF NOT EXISTS freeze_blocks_period INT;

ALTER TABLE parties
    DROP COLUMN IF EXISTS verified;
ALTER TABLE parties
    ADD COLUMN IF NOT EXISTS status INT;
ALTER TABLE parties
    ADD COLUMN IF NOT EXISTS violations_count INT;
ALTER TABLE parties
    ADD COLUMN IF NOT EXISTS freeze_end_block INT;
ALTER TABLE parties
    ADD COLUMN IF NOT EXISTS delegator TEXT;

CREATE TABLE violation_report
(
    index          TEXT UNIQUE PRIMARY KEY NOT NULL,
    session_id     TEXT                    NOT NULL,
    offender       TEXT                    NOT NULL REFERENCES account (address),
    sender         TEXT                    NOT NULL REFERENCES account (address),
    violation_type INT                     NOT NULL,
    msg            TEXT                    NOT NULL DEFAULT ''
);

-- +migrate Down
DROP TABLE IF EXISTS violation_report;

ALTER TABLE parties
    DROP COLUMN IF EXISTS delegator;
ALTER TABLE parties
    DROP COLUMN IF EXISTS freeze_end_block;
ALTER TABLE parties
    DROP COLUMN IF EXISTS violations_count;
ALTER TABLE parties
    DROP COLUMN IF EXISTS status;
ALTER TABLE parties
    ADD COLUMN IF NOT EXISTS verified BOOLEAN;

ALTER TABLE rarimocore_params
    DROP COLUMN IF EXISTS freezeBlocksPeriod;
ALTER TABLE rarimocore_params
    DROP COLUMN IF EXISTS maxViolationsCount;
ALTER TABLE rarimocore_params
    DROP COLUMN IF EXISTS stakeDenom;
ALTER TABLE rarimocore_params
    DROP COLUMN IF EXISTS stakeAmount;

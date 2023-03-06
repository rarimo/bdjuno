-- +migrate Up

ALTER TABLE rarimocore_params
    ADD COLUMN IF NOT EXISTS vote_threshold TEXT;
ALTER TABLE rarimocore_params
    ADD COLUMN IF NOT EXISTS vote_quorum TEXT;

-- +migrate Down
ALTER TABLE rarimocore_params
    DROP COLUMN IF EXISTS vote_threshold;
ALTER TABLE rarimocore_params
    DROP COLUMN IF EXISTS vote_quorum;

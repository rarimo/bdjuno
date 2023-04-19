package types

import (
	"encoding/json"
)

// MultisigProposalRow represents a single row of the "multisig_proposal" table
type MultisigProposalRow struct {
	Id               uint64          `db:"id"`
	Proposer         string          `db:"proposer"`
	Group            string          `db:"group"`
	SubmitBlock      uint64          `db:"submit_block"`
	VotingEndBlock   uint64          `db:"voting_end_block"`
	Status           int32           `db:"status"`
	FinalTallyResult json.RawMessage `db:"final_tally_result"`
	Messages         json.RawMessage `db:"messages"`
}

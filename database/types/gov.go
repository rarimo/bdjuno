package types

// GovParamsRow represents a single row of the "gov_params" table
type GovParamsRow struct {
	OneRowID      bool   `db:"one_row_id"`
	DepositParams string `db:"deposit_params"`
	VotingParams  string `db:"voting_params"`
	TallyParams   string `db:"tally_params"`
	Height        int64  `db:"height"`
}

// --------------------------------------------------------------------------------------------------------------------

// ProposalRow represents a single row inside the proposal table
type ProposalRow struct {
	Content          string `db:"content"`
	ProposalID       uint64 `db:"id"`
	SubmitBlock      uint64 `db:"submit_block"`
	DepositEndBlock  uint64 `db:"deposit_end_block"`
	VotingStartBlock uint64 `db:"voting_start_block"`
	VotingEndBlock   uint64 `db:"voting_end_block"`
	Proposer         string `db:"proposer_address"`
	Status           string `db:"status"`
	Metadata         string `db:"metadata"`
}

// NewProposalRow allows to easily create a new ProposalRow
func NewProposalRow(
	proposalID uint64,
	content string,
	submitBlock,
	depositBlock,
	votingStartBlock,
	votingEndBlock uint64,
	proposer,
	status,
	metadata string,
) ProposalRow {
	return ProposalRow{
		Content:          content,
		ProposalID:       proposalID,
		SubmitBlock:      submitBlock,
		DepositEndBlock:  depositBlock,
		VotingStartBlock: votingStartBlock,
		VotingEndBlock:   votingEndBlock,
		Proposer:         proposer,
		Status:           status,
		Metadata:         metadata,
	}
}

// Equals return true if two ProposalRow are the same
func (w ProposalRow) Equals(v ProposalRow) bool {
	return w.ProposalID == v.ProposalID &&
		w.SubmitBlock == v.SubmitBlock &&
		w.DepositEndBlock == v.DepositEndBlock &&
		w.VotingStartBlock == v.VotingStartBlock &&
		w.VotingEndBlock == v.VotingEndBlock &&
		w.Proposer == v.Proposer &&
		w.Status == v.Status &&
		w.Metadata == v.Metadata

}

// TallyResultRow represents a single row inside the tally_result table
type TallyResultRow struct {
	ProposalID int64  `db:"proposal_id"`
	Yes        string `db:"yes"`
	Abstain    string `db:"abstain"`
	No         string `db:"no"`
	NoWithVeto string `db:"no_with_veto"`
	Height     int64  `db:"height"`
}

// NewTallyResultRow return a new TallyResultRow instance
func NewTallyResultRow(
	proposalID int64,
	yes string,
	abstain string,
	no string,
	noWithVeto string,
	height int64,
) TallyResultRow {
	return TallyResultRow{
		ProposalID: proposalID,
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
		Height:     height,
	}
}

// Equals return true if two TallyResultRow are the same
func (w TallyResultRow) Equals(v TallyResultRow) bool {
	return w.ProposalID == v.ProposalID &&
		w.Yes == v.Yes &&
		w.Abstain == v.Abstain &&
		w.No == v.No &&
		w.NoWithVeto == v.NoWithVeto &&
		w.Height == v.Height
}

// VoteRow represents a single row inside the vote table
type VoteRow struct {
	ProposalID int64  `db:"proposal_id"`
	Voter      string `db:"voter_address"`
	Option     string `db:"option"`
	Height     int64  `db:"height"`
}

// NewVoteRow allows to easily create a new VoteRow
func NewVoteRow(
	proposalID int64,
	voter string,
	option string,
	height int64,
) VoteRow {
	return VoteRow{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
		Height:     height,
	}
}

// Equals return true if two VoteRow are the same
func (w VoteRow) Equals(v VoteRow) bool {
	return w.ProposalID == v.ProposalID &&
		w.Voter == v.Voter &&
		w.Option == v.Option &&
		w.Height == v.Height
}

// DepositRow represents a single row inside the deposit table
type DepositRow struct {
	ProposalID int64   `db:"proposal_id"`
	Depositor  string  `db:"depositor_address"`
	Amount     DbCoins `db:"amount"`
	Height     int64   `db:"height"`
}

// NewDepositRow allows to easily create a new NewDepositRow
func NewDepositRow(
	proposalID int64,
	depositor string,
	amount DbCoins,
	height int64,
) DepositRow {
	return DepositRow{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
		Height:     height,
	}
}

// Equals return true if two VoteDepositRow are the same
func (w DepositRow) Equals(v DepositRow) bool {
	return w.ProposalID == v.ProposalID &&
		w.Depositor == v.Depositor &&
		w.Amount.Equal(&v.Amount) &&
		w.Height == v.Height
}

// --------------------------------------------------------------------------------------------------------------------

type ProposalStakingPoolSnapshotRow struct {
	ProposalID      uint64 `db:"proposal_id"`
	BondedTokens    int64  `db:"bonded_tokens"`
	NotBondedTokens int64  `db:"not_bonded_tokens"`
	Height          int64  `db:"height"`
}

func NewProposalStakingPoolSnapshotRow(proposalID uint64, bondedTokens, notBondedTokens, height int64) ProposalStakingPoolSnapshotRow {
	return ProposalStakingPoolSnapshotRow{
		ProposalID:      proposalID,
		BondedTokens:    bondedTokens,
		NotBondedTokens: notBondedTokens,
		Height:          height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type ProposalValidatorVotingPowerSnapshotRow struct {
	ID               int64  `db:"id"`
	ProposalID       int64  `db:"proposal_id"`
	ValidatorAddress string `db:"validator_address"`
	VotingPower      int64  `db:"voting_power"`
	Status           int    `db:"status"`
	Jailed           bool   `db:"jailed"`
	Height           int64  `db:"height"`
}

func NewProposalValidatorVotingPowerSnapshotRow(
	id int64, proposalID int64, validatorAddr string, votingPower int64, status int, jailed bool, height int64,
) ProposalValidatorVotingPowerSnapshotRow {
	return ProposalValidatorVotingPowerSnapshotRow{
		ID:               id,
		ProposalID:       proposalID,
		ValidatorAddress: validatorAddr,
		VotingPower:      votingPower,
		Status:           status,
		Jailed:           jailed,
		Height:           height,
	}
}

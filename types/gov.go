package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalStatusInvalid = "PROPOSAL_STATUS_INVALID"
)

// DepositParams contains the data of the deposit parameters of the x/gov module
type DepositParams struct {
	MinDeposit       sdk.Coins `json:"min_deposit,omitempty" yaml:"min_deposit"`
	MaxDepositPeriod uint64    `json:"max_deposit_period,omitempty" yaml:"max_deposit_period"`
}

// NewDepositParam allows to build a new DepositParams
func NewDepositParam(d govtypes.DepositParams) DepositParams {
	return DepositParams{
		MinDeposit:       d.MinDeposit,
		MaxDepositPeriod: d.MaxDepositPeriod,
	}
}

// VotingParams contains the voting parameters of the x/gov module
type VotingParams struct {
	VotingPeriod uint64 `json:"voting_period,omitempty" yaml:"voting_period"`
}

// NewVotingParams allows to build a new VotingParams instance
func NewVotingParams(v govtypes.VotingParams) VotingParams {
	return VotingParams{
		VotingPeriod: v.VotingPeriod,
	}
}

// GovParams contains the data of the x/gov module parameters
type GovParams struct {
	DepositParams DepositParams `json:"deposit_params" yaml:"deposit_params"`
	VotingParams  VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   TallyParams   `json:"tally_params" yaml:"tally_params"`
	Height        int64         `json:"height" ymal:"height"`
}

// TallyParams contains the tally parameters of the x/gov module
type TallyParams struct {
	Quorum        sdk.Dec `json:"quorum,omitempty"`
	Threshold     sdk.Dec `json:"threshold,omitempty"`
	VetoThreshold sdk.Dec `json:"veto_threshold,omitempty" yaml:"veto_threshold"`
}

// NewTallyParams allows to build a new TallyParams instance
func NewTallyParams(t govtypes.TallyParams) TallyParams {
	return TallyParams{
		Quorum:        t.Quorum,
		Threshold:     t.Threshold,
		VetoThreshold: t.VetoThreshold,
	}
}

// NewGovParams allows to build a new GovParams instance
func NewGovParams(votingParams VotingParams, depositParams DepositParams, tallyParams TallyParams, height int64) *GovParams {
	return &GovParams{
		DepositParams: depositParams,
		VotingParams:  votingParams,
		TallyParams:   tallyParams,
		Height:        height,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// Proposal represents a single governance proposal
type Proposal struct {
	ProposalRoute    string
	ProposalType     string
	ProposalID       uint64
	Content          govtypes.Content
	Status           string
	SubmitBlock      uint64
	DepositEndBlock  uint64
	VotingStartBlock uint64
	VotingEndBlock   uint64
	Proposer         string
}

// NewProposal return a new Proposal instance
func NewProposal(
	proposalID uint64,
	proposalRoute string,
	proposalType string,
	content govtypes.Content,
	status string,
	submitBlock,
	depositEndBlock,
	votingStartBlock,
	votingEndBlock uint64,
	proposer string,
) Proposal {
	return Proposal{
		Content:          content,
		ProposalRoute:    proposalRoute,
		ProposalType:     proposalType,
		ProposalID:       proposalID,
		Status:           status,
		SubmitBlock:      submitBlock,
		DepositEndBlock:  depositEndBlock,
		VotingStartBlock: votingStartBlock,
		VotingEndBlock:   votingEndBlock,
		Proposer:         proposer,
	}
}

// Equal tells whether p and other contain the same data
func (p Proposal) Equal(other Proposal) bool {
	return p.ProposalRoute == other.ProposalRoute &&
		p.ProposalType == other.ProposalType &&
		p.ProposalID == other.ProposalID &&
		p.Content.String() == other.Content.String() &&
		p.Status == other.Status &&
		p.SubmitBlock == other.SubmitBlock &&
		p.DepositEndBlock == other.DepositEndBlock &&
		p.VotingStartBlock == other.VotingStartBlock &&
		p.VotingEndBlock == other.VotingEndBlock &&
		p.Proposer == other.Proposer
}

// ProposalUpdate contains the data that should be used when updating a governance proposal
type ProposalUpdate struct {
	ProposalID       uint64
	Status           string
	VotingStartBlock uint64
	VotingEndBlock   uint64
}

// NewProposalUpdate allows to build a new ProposalUpdate instance
func NewProposalUpdate(
	proposalID uint64, status string, votingStartBlock, votingEndBlock uint64,
) ProposalUpdate {
	return ProposalUpdate{
		ProposalID:       proposalID,
		Status:           status,
		VotingStartBlock: votingStartBlock,
		VotingEndBlock:   votingEndBlock,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// Deposit contains the data of a single deposit made towards a proposal
type Deposit struct {
	ProposalID uint64
	Depositor  string
	Amount     sdk.Coins
	Height     int64
}

// NewDeposit return a new Deposit instance
func NewDeposit(
	proposalID uint64,
	depositor string,
	amount sdk.Coins,
	height int64,
) Deposit {
	return Deposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
		Height:     height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// Vote contains the data of a single proposal vote
type Vote struct {
	ProposalID uint64
	Voter      string
	Option     govtypes.VoteOption
	Height     int64
}

// NewVote return a new Vote instance
func NewVote(
	proposalID uint64,
	voter string,
	option govtypes.VoteOption,
	height int64,
) Vote {
	return Vote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
		Height:     height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// TallyResult contains the data about the final results of a proposal
type TallyResult struct {
	ProposalID uint64
	Yes        string
	Abstain    string
	No         string
	NoWithVeto string
	Height     int64
}

// NewTallyResult return a new TallyResult instance
func NewTallyResult(
	proposalID uint64,
	yes string,
	abstain string,
	no string,
	noWithVeto string,
	height int64,
) TallyResult {
	return TallyResult{
		ProposalID: proposalID,
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
		Height:     height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// ProposalStakingPoolSnapshot contains the data about a single staking pool snapshot to be associated with a proposal
type ProposalStakingPoolSnapshot struct {
	ProposalID uint64
	Pool       *Pool
}

// NewProposalStakingPoolSnapshot returns a new ProposalStakingPoolSnapshot instance
func NewProposalStakingPoolSnapshot(proposalID uint64, pool *Pool) ProposalStakingPoolSnapshot {
	return ProposalStakingPoolSnapshot{
		ProposalID: proposalID,
		Pool:       pool,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// ProposalValidatorStatusSnapshot represents a single snapshot of the status of a validator associated
// with a single proposal
type ProposalValidatorStatusSnapshot struct {
	ProposalID           uint64
	ValidatorConsAddress string
	ValidatorVotingPower int64
	ValidatorStatus      int
	ValidatorJailed      bool
	Height               int64
}

// NewProposalValidatorStatusSnapshot returns a new ProposalValidatorStatusSnapshot instance
func NewProposalValidatorStatusSnapshot(
	proposalID uint64,
	validatorConsAddr string,
	validatorVotingPower int64,
	validatorStatus int,
	validatorJailed bool,
	height int64,
) ProposalValidatorStatusSnapshot {
	return ProposalValidatorStatusSnapshot{
		ProposalID:           proposalID,
		ValidatorStatus:      validatorStatus,
		ValidatorConsAddress: validatorConsAddr,
		ValidatorVotingPower: validatorVotingPower,
		ValidatorJailed:      validatorJailed,
		Height:               height,
	}
}

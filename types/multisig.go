package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	multisigtypes "github.com/rarimo/rarimo-core/x/multisig/types"
)

//--------------------------------------------------------

// MultisigParams contains the data of the x/multisig module params instance
type MultisigParams struct {
	GroupSequence    uint64 `json:"group_sequence,omitempty" yaml:"group_sequence,omitempty"`
	ProposalSequence uint64 `json:"proposal_sequence,omitempty" yaml:"proposal_sequence,omitempty"`
	PrunePeriod      uint64 `json:"prune_period,omitempty" yaml:"prune_period,omitempty"`
	VotingPeriod     uint64 `json:"voting_period,omitempty" yaml:"voting_period,omitempty"`
	Height           int64  `json:"height,omitempty" yaml:"height,omitempty"`
}

// MultisigParamsFromCore allows to build a new MultisigParams instance from an oraclemanagertypes.Params instance
func MultisigParamsFromCore(p multisigtypes.Params, height int64) *MultisigParams {
	return &MultisigParams{
		GroupSequence:    p.GroupSequence,
		ProposalSequence: p.ProposalSequence,
		PrunePeriod:      p.PrunePeriod,
		VotingPeriod:     p.VotingPeriod,
		Height:           height,
	}
}

//--------------------------------------------------------

type Group struct {
	Account   string   `json:"account,omitempty" yaml:"account,omitempty"`
	Members   []string `json:"members,omitempty" yaml:"members,omitempty"`
	Threshold uint64   `json:"threshold,omitempty" yaml:"threshold,omitempty"`
}

// GroupFromCore allows to build a new Group instance from a multisigtypes.Group instance
func GroupFromCore(g multisigtypes.Group) *Group {
	return &Group{
		Account:   g.Account,
		Members:   g.Members,
		Threshold: g.Threshold,
	}
}

//--------------------------------------------------------

type MultisigProposal struct {
	Id               uint64                       `json:"id,omitempty" yaml:"id,omitempty"`
	Proposer         string                       `json:"proposer,omitempty" yaml:"proposer,omitempty"`
	Group            string                       `json:"group,omitempty" yaml:"group,omitempty"`
	SubmitBlock      uint64                       `json:"submit_block,omitempty" yaml:"submit_block,omitempty"`
	VotingEndBlock   uint64                       `json:"voting_end_block,omitempty" yaml:"voting_end_block,omitempty"`
	Status           multisigtypes.ProposalStatus `json:"status,omitempty" yaml:"status,omitempty"`
	FinalTallyResult *multisigtypes.TallyResult   `json:"final_tally_result,omitempty" yaml:"final_tally_result,omitempty"`
	Messages         []*types.Any                 `json:"messages,omitempty" yaml:"messages,omitempty"`
}

// MultisigProposalFromCore allows to build a new MultisigProposal instance from a multisigtypes.Proposal instance
func NewMultisigProposal(id uint64, proposer, group string, submitBlock, votingEndBlock uint64, status int32, result, messages json.RawMessage) (*MultisigProposal, error) {
	p := &MultisigProposal{
		Id:             id,
		Proposer:       proposer,
		Group:          group,
		SubmitBlock:    submitBlock,
		VotingEndBlock: votingEndBlock,
		Status:         multisigtypes.ProposalStatus(status),
	}

	if result != nil {
		var tally multisigtypes.TallyResult
		if err := json.Unmarshal(result, &tally); err != nil {
			return nil, err
		}
		p.FinalTallyResult = &tally
	}

	return p, nil
}

func MultisigProposalFromCore(p multisigtypes.Proposal) *MultisigProposal {
	return &MultisigProposal{
		Id:               p.Id,
		Proposer:         p.Proposer,
		Group:            p.Group,
		SubmitBlock:      p.SubmitBlock,
		VotingEndBlock:   p.VotingEndBlock,
		Status:           p.Status,
		FinalTallyResult: p.FinalTallyResult,
		Messages:         p.Messages,
	}
}

//--------------------------------------------------------

type MultisigProposalVote struct {
	Index       string                   `json:"index,omitempty" yaml:"index,omitempty"`
	ProposalId  uint64                   `json:"proposal_id,omitempty" yaml:"proposal_id,omitempty"`
	Voter       string                   `json:"voter,omitempty" yaml:"voter,omitempty"`
	Option      multisigtypes.VoteOption `json:"option,omitempty" yaml:"option,omitempty"`
	SubmitBlock uint64                   `json:"submit_block,omitempty" yaml:"submit_block,omitempty"`
}

// MultisigProposalVoteFromCore allows to build a new MultisigProposalVote instance from a multisigtypes.Vote instance
func MultisigProposalVoteFromCore(p multisigtypes.Vote) *MultisigProposalVote {
	return &MultisigProposalVote{
		Index:      hexutil.Encode(multisigtypes.VoteKey(p.ProposalId, p.Voter)),
		ProposalId: p.ProposalId,
		Voter:      p.Voter,
		Option:     p.Option,
	}
}

//--------------------------------------------------------

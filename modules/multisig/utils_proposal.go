package multisig

import (
	"fmt"
	"gitlab.com/rarimo/bdjuno/types"
	multisigtypes "gitlab.com/rarimo/rarimo-core/x/multisig/types"
)

func (m *Module) saveProposals(slice []multisigtypes.Proposal) error {
	// Save the proposals
	proposals := make([]*types.MultisigProposal, len(slice))
	for i, proposal := range slice {
		proposals[i] = types.MultisigProposalFromCore(proposal)
	}

	err := m.db.SaveMultisigProposals(proposals)
	if err != nil {
		return fmt.Errorf("error while storing genesis multisig proposals: %s", err)
	}

	return nil
}

func (m *Module) saveProposal(height int64, id uint64) error {
	proposal, err := m.source.Proposal(height, id)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	return m.saveProposals([]multisigtypes.Proposal{proposal})
}

func (m *Module) saveVotes(slice []multisigtypes.Vote) error {
	// Save the votes
	votes := make([]*types.MultisigProposalVote, len(slice))
	for i, vote := range slice {
		votes[i] = types.MultisigProposalVoteFromCore(vote)
	}

	err := m.db.SaveMultisigProposalVotes(votes)
	if err != nil {
		return fmt.Errorf("error while storing genesis multisig votes: %s", err)
	}

	return nil
}

func (m *Module) saveVote(height int64, proposalId uint64, voter string) error {
	vote, err := m.source.Vote(height, proposalId, voter)
	if err != nil {
		return fmt.Errorf("error while getting proposal vote: %s", err)
	}

	return m.saveVotes([]multisigtypes.Vote{vote})
}

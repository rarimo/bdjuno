package source

import multisigtypes "gitlab.com/rarimo/rarimo-core/x/multisig/types"

type Source interface {
	Params(height int64) (multisigtypes.Params, error)
	Group(height int64, account string) (multisigtypes.Group, error)
	Proposal(height int64, id uint64) (multisigtypes.Proposal, error)
	ProposalAll(height int64) ([]multisigtypes.Proposal, error)
	Vote(height int64, proposalId uint64, voter string) (multisigtypes.Vote, error)
}

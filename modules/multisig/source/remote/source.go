package remote

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v4/node/remote"
	multisigsource "gitlab.com/rarimo/bdjuno/modules/multisig/source"
	multisigtypes "gitlab.com/rarimo/rarimo-core/x/multisig/types"
)

var (
	_ multisigsource.Source = &Source{}
)

// Source implements multisigsource.Source using a remote node
type Source struct {
	*remote.Source
	client multisigtypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, client multisigtypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// Params implements multisigsource.Source
func (s Source) Params(height int64) (multisigtypes.Params, error) {
	res, err := s.client.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&multisigtypes.QueryParamsRequest{},
	)
	if err != nil {
		return multisigtypes.Params{}, err
	}

	return res.Params, err
}

// Group implements multisigsource.Source
func (s Source) Group(height int64, account string) (multisigtypes.Group, error) {
	res, err := s.client.Group(
		remote.GetHeightRequestContext(s.Ctx, height),
		&multisigtypes.QueryGetGroupRequest{Account: account},
	)
	if err != nil {
		return multisigtypes.Group{}, err
	}

	return res.Group, err
}

// Proposal implements multisigsource.Source
func (s Source) Proposal(height int64, id uint64) (multisigtypes.Proposal, error) {
	res, err := s.client.Proposal(
		remote.GetHeightRequestContext(s.Ctx, height),
		&multisigtypes.QueryGetProposalRequest{ProposalId: id},
	)
	if err != nil {
		return multisigtypes.Proposal{}, err
	}

	return res.Proposal, err
}

// ProposalAll implements multisigsource.Source
func (s Source) ProposalAll(height int64) ([]multisigtypes.Proposal, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var proposals []multisigtypes.Proposal
	var nextKey []byte
	var stop = false

	for !stop {
		res, err := s.client.ProposalAll(
			ctx,
			&multisigtypes.QueryAllProposalRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting proposal all: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		proposals = append(proposals, res.Proposal...)
	}

	return proposals, nil
}

// Vote implements multisigsource.Source
func (s Source) Vote(height int64, proposalId uint64, voter string) (multisigtypes.Vote, error) {
	res, err := s.client.Vote(
		remote.GetHeightRequestContext(s.Ctx, height),
		&multisigtypes.QueryGetVoteRequest{ProposalId: proposalId, Voter: voter},
	)
	if err != nil {
		return multisigtypes.Vote{}, err
	}

	return res.Vote, err
}

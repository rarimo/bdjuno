package remote

import (
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/forbole/juno/v4/node/remote"

	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	govsource "github.com/rarimo/bdjuno/modules/gov/source"
)

var (
	_ govsource.Source = &Source{}
)

// Source implements govsource.Source using a remote node
type Source struct {
	*remote.Source
	q        govtypesv1.QueryClient
	qv1beta1 govtypesv1beta1.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, govKeeper govtypesv1.QueryClient, govKeeperv1beta1 govtypesv1beta1.QueryClient) *Source {
	return &Source{
		Source:   source,
		q:        govKeeper,
		qv1beta1: govKeeperv1beta1,
	}
}

// Proposal implements govsource.Source
func (s Source) Proposal(height int64, id uint64) (govtypesv1beta1.Proposal, error) {
	res, err := s.qv1beta1.Proposal(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1beta1.QueryProposalRequest{ProposalId: id},
	)
	if err != nil {
		return govtypesv1beta1.Proposal{}, err
	}

	return res.Proposal, err
}

func (s Source) ProposalV1(height int64, id uint64) (govtypesv1.Proposal, error) {
	res, err := s.q.Proposal(remote.GetHeightRequestContext(s.Ctx, height), &govtypesv1.QueryProposalRequest{ProposalId: id})
	if err != nil {
		return govtypesv1.Proposal{}, err
	}

	return *res.Proposal, err
}

// ProposalDeposit implements govsource.Source
func (s Source) ProposalDeposit(height int64, id uint64, depositor string) (*govtypesv1.Deposit, error) {
	res, err := s.q.Deposit(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryDepositRequest{ProposalId: id, Depositor: depositor},
	)
	if err != nil {
		return nil, err
	}

	return res.Deposit, nil
}

// TallyResult implements govsource.Source
func (s Source) TallyResult(height int64, proposalID uint64) (*govtypesv1.TallyResult, error) {
	res, err := s.q.TallyResult(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryTallyResultRequest{ProposalId: proposalID},
	)
	if err != nil {
		return nil, err
	}

	return res.Tally, nil
}

// DepositParams implements govsource.Source
func (s Source) DepositParams(height int64) (*govtypesv1.DepositParams, error) {
	res, err := s.q.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryParamsRequest{ParamsType: govtypesv1.ParamDeposit},
	)
	if err != nil {
		return nil, err
	}

	return res.DepositParams, nil
}

// VotingParams implements govsource.Source
func (s Source) VotingParams(height int64) (*govtypesv1.VotingParams, error) {
	res, err := s.q.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryParamsRequest{ParamsType: govtypesv1.ParamVoting},
	)
	if err != nil {
		return nil, err
	}

	return res.VotingParams, nil
}

// TallyParams implements govsource.Source
func (s Source) TallyParams(height int64) (*govtypesv1.TallyParams, error) {
	res, err := s.q.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryParamsRequest{ParamsType: govtypesv1.ParamTallying},
	)
	if err != nil {
		return nil, err
	}

	return res.TallyParams, nil
}

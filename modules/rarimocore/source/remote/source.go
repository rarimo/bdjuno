package remote

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"
	rarimocoresource "gitlab.com/rarimo/bdjuno/modules/rarimocore/source"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

var (
	_ rarimocoresource.Source = &Source{}
)

// Source implements rarimocoresource.Source using a remote node
type Source struct {
	*remote.Source
	rarimocoreClient rarimocoretypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, rarimocoreClient rarimocoretypes.QueryClient) *Source {
	return &Source{
		Source:           source,
		rarimocoreClient: rarimocoreClient,
	}
}

// Params implements rarimocoresource.Source
func (s Source) Params(height int64) (rarimocoretypes.Params, error) {
	res, err := s.rarimocoreClient.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&rarimocoretypes.QueryParamsRequest{},
	)
	if err != nil {
		return rarimocoretypes.Params{}, err
	}

	return res.Params, err
}

// Operation implements rarimocoresource.Source
func (s Source) Operation(height int64, index string) (rarimocoretypes.Operation, error) {
	res, err := s.rarimocoreClient.Operation(
		remote.GetHeightRequestContext(s.Ctx, height),
		&rarimocoretypes.QueryGetOperationRequest{Index: index},
	)
	if err != nil {
		return rarimocoretypes.Operation{}, err
	}

	return res.Operation, err
}

// OperationAll implements rarimocoresource.Source
func (s Source) OperationAll(height int64) ([]rarimocoretypes.Operation, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var operations []rarimocoretypes.Operation
	var nextKey []byte
	var stop = false

	for !stop {
		res, err := s.rarimocoreClient.OperationAll(
			ctx,
			&rarimocoretypes.QueryAllOperationRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting operation all: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		operations = append(operations, res.Operation...)
	}

	return operations, nil
}

// Confirmation implements rarimocoresource.Source
func (s Source) Confirmation(height int64, root string) (rarimocoretypes.Confirmation, error) {
	res, err := s.rarimocoreClient.Confirmation(
		remote.GetHeightRequestContext(s.Ctx, height),
		&rarimocoretypes.QueryGetConfirmationRequest{Root: root},
	)
	if err != nil {
		return rarimocoretypes.Confirmation{}, err
	}

	return res.Confirmation, err
}

// ConfirmationAll implements rarimocoresource.Source
func (s Source) ConfirmationAll(height int64) ([]rarimocoretypes.Confirmation, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var confirmations []rarimocoretypes.Confirmation
	var nextKey []byte
	var stop = false

	for !stop {
		res, err := s.rarimocoreClient.ConfirmationAll(
			ctx,
			&rarimocoretypes.QueryAllConfirmationRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting confirmation all: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		confirmations = append(confirmations, res.Confirmation...)
	}

	return confirmations, nil
}

package remote

import (
	"github.com/forbole/juno/v4/node/remote"
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

// ViolationReport implements rarimocoresource.Source
func (s Source) ViolationReport(height int64, sessionId, offender, sender string, violationType rarimocoretypes.ViolationType) (rarimocoretypes.ViolationReport, error) {
	res, err := s.rarimocoreClient.ViolationReport(
		remote.GetHeightRequestContext(s.Ctx, height),
		&rarimocoretypes.QueryGetViolationReportRequest{
			SessionId:     sessionId,
			Offender:      offender,
			ViolationType: violationType,
			Sender:        sender,
		},
	)
	if err != nil {
		return rarimocoretypes.ViolationReport{}, err
	}

	return res.ViolationReport, err
}

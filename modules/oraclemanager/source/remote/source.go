package remote

import (
	"github.com/forbole/juno/v4/node/remote"
	oraclemanagersource "gitlab.com/rarimo/bdjuno/modules/oraclemanager/source"
	oraclemanagertypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
)

var (
	_ oraclemanagersource.Source = &Source{}
)

// Source implements oraclemanagersource.Source using a remote node
type Source struct {
	*remote.Source
	client oraclemanagertypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, client oraclemanagertypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// Params implements oraclemanagersource.Source
func (s Source) Params(height int64) (oraclemanagertypes.Params, error) {
	res, err := s.client.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&oraclemanagertypes.QueryParamsRequest{},
	)
	if err != nil {
		return oraclemanagertypes.Params{}, err
	}

	return res.Params, err
}

// Oracle implements oraclemanagersource.Source
func (s Source) Oracle(height int64, chain, account string) (oraclemanagertypes.Oracle, error) {
	res, err := s.client.Oracle(
		remote.GetHeightRequestContext(s.Ctx, height),
		&oraclemanagertypes.QueryGetOracleRequest{
			Chain:   chain,
			Address: account,
		},
	)
	if err != nil {
		return oraclemanagertypes.Oracle{}, err
	}

	return res.Oracle, err
}

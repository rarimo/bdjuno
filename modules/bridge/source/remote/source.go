package remote

import (
	"github.com/forbole/juno/v4/node/remote"
	bridgesource "gitlab.com/rarimo/bdjuno/modules/bridge/source"
	bridgetypes "gitlab.com/rarimo/rarimo-core/x/bridge/types"
)

var (
	_ bridgesource.Source = &Source{}
)

// Source implements bridgesource.Source using a remote node
type Source struct {
	*remote.Source
	client bridgetypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, client bridgetypes.QueryClient) *Source {
	return &Source{
		Source: source,
		client: client,
	}
}

// Params implements bridgesource.Source
func (s Source) Params(height int64) (bridgetypes.Params, error) {
	res, err := s.client.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&bridgetypes.QueryParamsRequest{},
	)
	if err != nil {
		return bridgetypes.Params{}, err
	}

	return res.Params, err
}

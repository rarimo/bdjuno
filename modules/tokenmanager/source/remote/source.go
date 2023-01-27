package remote

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"
	tokenmanagersource "gitlab.com/rarimo/bdjuno/modules/tokenmanager/source"
	"gitlab.com/rarimo/bdjuno/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

var (
	_ tokenmanagersource.Source = &Source{}
)

// Source implements tokenmanagersource.Source using a remote node
type Source struct {
	*remote.Source
	tokenmanagerClient tokenmanagertypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, tokenmanagerClient tokenmanagertypes.QueryClient) *Source {
	return &Source{
		Source:             source,
		tokenmanagerClient: tokenmanagerClient,
	}
}

// Params implements tokenmanagersource.Source
func (s Source) Params(height int64) (tokenmanagertypes.Params, error) {
	res, err := s.tokenmanagerClient.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&tokenmanagertypes.QueryParamsRequest{},
	)
	if err != nil {
		return tokenmanagertypes.Params{}, err
	}

	return res.Params, err
}

// Item implements tokenmanagersource.Source
func (s Source) Item(height int64, index string) (tokenmanagertypes.Item, error) {
	res, err := s.tokenmanagerClient.Item(
		remote.GetHeightRequestContext(s.Ctx, height),
		&tokenmanagertypes.QueryGetItemRequest{
			Index: index,
		},
	)
	if err != nil {
		return tokenmanagertypes.Item{}, err
	}

	return res.Item, err
}

// ItemAll implements tokenmanagersource.Source
func (s Source) ItemAll(height int64) ([]tokenmanagertypes.Item, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var items []tokenmanagertypes.Item
	var nextKey []byte
	var stop = false

	for !stop {
		res, err := s.tokenmanagerClient.ItemAll(
			ctx,
			&tokenmanagertypes.QueryAllItemRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting item all: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		items = append(items, res.Item...)
	}

	return items, nil
}

// Collection implements tokenmanagersource.Source
func (s Source) Collection(height int64, index string) (tokenmanagertypes.Collection, error) {
	res, err := s.tokenmanagerClient.Collection(
		remote.GetHeightRequestContext(s.Ctx, height),
		&tokenmanagertypes.QueryGetCollectionRequest{Index: index},
	)
	if err != nil {
		return tokenmanagertypes.Collection{}, err
	}

	return res.Collection, err
}

// CollectionAll implements tokenmanagersource.Source
func (s Source) CollectionAll(height int64) ([]tokenmanagertypes.Collection, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var collections []tokenmanagertypes.Collection
	var nextKey []byte
	var stop = false

	for !stop {
		res, err := s.tokenmanagerClient.CollectionAll(
			ctx,
			&tokenmanagertypes.QueryAllCollectionRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting collection all: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		collections = append(collections, res.Collection...)
	}

	return collections, nil
}

// CollectionData implements tokenmanagersource.Source
func (s Source) CollectionData(height int64, index types.CollectionDataIndex) (tokenmanagertypes.CollectionData, error) {
	res, err := s.tokenmanagerClient.CollectionData(
		remote.GetHeightRequestContext(s.Ctx, height),
		&tokenmanagertypes.QueryGetCollectionDataRequest{
			Chain:   index.Chain,
			Address: index.Address,
		},
	)
	if err != nil {
		return tokenmanagertypes.CollectionData{}, err
	}

	return res.Data, err
}

// CollectionDataAll implements tokenmanagersource.Source
func (s Source) CollectionDataAll(height int64) ([]tokenmanagertypes.CollectionData, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var collectionDatas []tokenmanagertypes.CollectionData
	var nextKey []byte
	var stop = false

	for !stop {
		res, err := s.tokenmanagerClient.CollectionDataAll(
			ctx,
			&tokenmanagertypes.QueryAllCollectionDataRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error while getting collection data all: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		collectionDatas = append(collectionDatas, res.Data...)
	}

	return collectionDatas, nil
}

package source

import (
	"github.com/rarimo/bdjuno/types"
	tokenmanagertypes "github.com/rarimo/rarimo-core/x/tokenmanager/types"
)

type Source interface {
	Params(height int64) (tokenmanagertypes.Params, error)
	GetNetwork(height int64, network string) (tokenmanagertypes.Network, bool)
	GetFeeToken(height int64, chain, contract string) (*tokenmanagertypes.FeeToken, error)
	Item(height int64, index string) (tokenmanagertypes.Item, error)
	ItemAll(height int64) ([]tokenmanagertypes.Item, error)
	OnChainItem(height int64, index types.OnChainItemIndex) (tokenmanagertypes.OnChainItem, error)
	Seed(height int64, seed string) (tokenmanagertypes.Seed, error)
	Collection(height int64, index string) (tokenmanagertypes.Collection, error)
	CollectionAll(height int64) ([]tokenmanagertypes.Collection, error)
	CollectionData(height int64, index types.CollectionDataIndex) (tokenmanagertypes.CollectionData, error)
	CollectionDataAll(height int64) ([]tokenmanagertypes.CollectionData, error)
}

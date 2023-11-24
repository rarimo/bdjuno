package tokenmanager

import (
	"fmt"
	"github.com/rarimo/bdjuno/types"
	tokenmanagertypes "github.com/rarimo/rarimo-core/x/tokenmanager/types"
)

func (m *Module) saveNetworks(params tokenmanagertypes.Params) (err error) {
	networks := make([]types.Network, len(params.Networks))

	for i, network := range params.Networks {
		net, err := types.NewNetwork(*network)
		if err != nil {
			return fmt.Errorf("error while building network: %s", err)
		}

		networks[i] = *net
	}

	return m.db.SaveNetworks(networks)
}

func (m *Module) saveCollections(collections []tokenmanagertypes.Collection) (err error) {
	list := make([]types.Collection, len(collections))

	for i, collection := range collections {
		list[i] = types.CollectionFromCore(collection)
	}

	return m.db.SaveCollections(list)
}

func (m *Module) saveCollectionDatas(datas []tokenmanagertypes.CollectionData) (err error) {
	list := make([]types.CollectionData, len(datas))

	for i, data := range datas {
		list[i] = types.CollectionDataFromCore(data)
	}

	return m.db.SaveCollectionDatas(list)
}

func (m *Module) saveItems(items []tokenmanagertypes.Item) (err error) {
	list := make([]types.Item, len(items))

	for i, item := range items {
		list[i] = types.ItemFromCore(item)
	}

	return m.db.SaveItems(list)
}

func (m *Module) saveOnChainItems(items []tokenmanagertypes.OnChainItem) (err error) {
	list := make([]types.OnChainItem, len(items))

	for i, item := range items {
		list[i] = types.OnChainItemFromCore(item)
	}

	return m.db.SaveOnChainItems(list)
}

func (m *Module) saveSeeds(seeds []tokenmanagertypes.Seed) (err error) {
	list := make([]types.Seed, len(seeds))

	for i, seed := range seeds {
		list[i] = types.SeedFromCore(seed)
	}

	return m.db.SaveSeeds(list)
}

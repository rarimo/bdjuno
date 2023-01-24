package tokenmanager

import (
	"gitlab.com/rarimo/bdjuno/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

func (m *Module) saveParams(params tokenmanagertypes.Params, height int64) (err error) {
	return m.db.SaveTokenManagerParams(types.NewTokenManagerParams(params, height))
}

func (m *Module) saveCollections(collections []tokenmanagertypes.Collection) (err error) {
	list := make([]types.Collection, len(collections))

	for i, collection := range collections {
		indexes := make([]*types.CollectionDataIndex, 0)

		for _, data := range collection.Data {
			if data == nil {
				continue
			}

			indexes = append(indexes, types.NewCollectionDataIndex(data.Chain, data.Address))
		}

		list[i] = types.NewCollection(
			collection.Index,
			types.NewCollectionMetadata(
				collection.Meta.Name,
				collection.Meta.Symbol,
				collection.Meta.MetadataURI,
			),
			indexes,
		)
	}

	return m.db.SaveCollections(list)
}

func (m *Module) saveCollectionDatas(datas []tokenmanagertypes.CollectionData) (err error) {
	list := make([]types.CollectionData, len(datas))

	for i, data := range datas {
		list[i] = types.NewCollectionData(
			types.NewCollectionDataIndex(data.Index.Chain, data.Index.Address),
			data.Collection,
			data.TokenType,
			data.Wrapped,
			data.Decimals,
		)
	}

	return m.db.SaveCollectionDatas(list)
}

func (m *Module) saveItems(items []tokenmanagertypes.Item) (err error) {
	list := make([]types.Item, len(items))

	for i, item := range items {
		params := make([]*types.ItemChainParams, 0)

		for _, param := range item.ChainParams {
			if param == nil {
				continue
			}

			params = append(params, types.NewItemChainParams(param.Chain, param.TokenID))
		}

		list[i] = types.NewItem(
			types.NewItemIndex(
				item.Index.Collection,
				item.Index.Name,
				item.Index.Symbol,
				item.Index.Uri,
			),
			types.NewItemMetadata(
				item.Meta.ImageUri,
				item.Meta.ImageHash,
				item.Meta.Seed,
			),
			params,
		)
	}

	return m.db.SaveItems(list)
}

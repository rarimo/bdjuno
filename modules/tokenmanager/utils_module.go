package tokenmanager

import (
	"fmt"
	"github.com/rarimo/bdjuno/types"
	tokenmanagertypes "github.com/rarimo/rarimo-core/x/tokenmanager/types"
)

func (m *Module) UpdateItems(items []*tokenmanagertypes.Item) (err error) {
	for _, item := range items {
		if item == nil {
			continue
		}

		err = m.updateItem(item)
		if err != nil {
			return fmt.Errorf("failed to update item in tokenmanager: %s", err)
		}
	}

	return nil
}

func (m *Module) updateItem(newItem *tokenmanagertypes.Item) error {
	oldItem, err := m.db.GetItem(newItem.Index)
	if err != nil {
		return fmt.Errorf("failed to get item in tokenmanager: %s", err)
	}

	if oldItem == nil {
		return fmt.Errorf("item not found in tokenmanager")
	}

	if oldItem.Meta.Seed != newItem.Meta.Seed {
		err = m.db.RemoveSeed(oldItem.Meta.Seed)
		if err != nil {
			return fmt.Errorf("failed to remove seed in tokenmanager: %s", err)
		}

		if newItem.Meta.Seed != "" {
			err = m.db.SaveSeeds([]types.Seed{types.NewSeed(newItem.Meta.Seed, newItem.Index)})
			if err != nil {
				return fmt.Errorf("failed to save seed in tokenmanager: %s", err)
			}
		}
	}

	err = m.db.UpsertItem(types.ItemFromCore(*newItem))
	if err != nil {
		return fmt.Errorf("failed to update item in tokenmanager: %s", err)
	}

	return nil
}

func (m *Module) RemoveItems(indexes []string) (err error) {
	for _, index := range indexes {
		err = m.removeItem(index)
		if err != nil {
			return fmt.Errorf("failed to remove item in tokenmanager: %s", err)
		}
	}

	return nil
}

func (m *Module) removeItem(index string) error {
	item, err := m.db.GetItem(index)
	if err != nil {
		return fmt.Errorf("failed to get item in tokenmanager: %s", err)
	}

	err = m.db.RemoveOnChainItems(index)
	if err != nil {
		return fmt.Errorf("failed to remove item in tokenmanager: %s", err)
	}

	if item.Meta.Seed != "" {
		err = m.db.RemoveSeed(item.Meta.Seed)
		if err != nil {
			return fmt.Errorf("failed to remove seed in tokenmanager: %s", err)
		}
	}

	err = m.db.RemoveItem(index)
	if err != nil {
		return fmt.Errorf("failed to remove item in tokenmanager: %s", err)
	}

	return nil
}

func (m *Module) CreateCollection(
	index string,
	meta *tokenmanagertypes.CollectionMetadata,
	data []*tokenmanagertypes.CollectionData,
	items []*tokenmanagertypes.Item,
	onChainItems []*tokenmanagertypes.OnChainItem,
) error {
	coreCollection := tokenmanagertypes.Collection{
		Index: index,
		Meta:  *meta,
		Data:  make([]*tokenmanagertypes.CollectionDataIndex, 0, len(data)),
	}

	for _, collectionData := range data {
		coreCollection.Data = append(coreCollection.Data, collectionData.Index)
	}

	collection := types.CollectionFromCore(coreCollection)

	datas := make([]tokenmanagertypes.CollectionData, 0, len(data))
	for _, collectionData := range data {
		datas = append(datas, *collectionData)
	}

	itemList := make([]tokenmanagertypes.Item, len(items))
	for i, item := range items {
		itemList[i] = *item
	}

	onChainItemsList := make([]tokenmanagertypes.OnChainItem, len(onChainItems))
	for i, onChainItem := range onChainItems {
		onChainItemsList[i] = *onChainItem
	}

	err := m.db.SaveCollections([]types.Collection{collection})
	if err != nil {
		return fmt.Errorf("failed to create collection: %s", err)
	}

	err = m.saveCollectionDatas(datas)
	if err != nil {
		return fmt.Errorf("failed to create collection datas: %s", err)
	}

	err = m.saveItems(itemList)
	if err != nil {
		return fmt.Errorf("failed to create items: %s", err)
	}

	err = m.saveOnChainItems(onChainItemsList)
	if err != nil {
		return fmt.Errorf("failed to create on chain items: %s", err)
	}

	return nil
}

func (m *Module) UpdateCollectionDatas(datas []*tokenmanagertypes.CollectionData) (err error) {
	for _, data := range datas {
		collectionData := types.CollectionDataFromCore(*data)
		err = m.db.UpdateCollectionData(collectionData)
		if err != nil {
			return fmt.Errorf("failed to update collection data: %s", err)
		}
	}

	return nil
}

func (m *Module) CreateCollectionDatas(height int64, datas []*tokenmanagertypes.CollectionData) (err error) {
	for _, data := range datas {
		err = m.createCollectionData(height, data)
		if err != nil {
			return fmt.Errorf("failed to create collection data: %s", err)
		}
	}

	return nil
}

func (m *Module) createCollectionData(height int64, data *tokenmanagertypes.CollectionData) error {
	err := m.db.SaveCollectionDatas([]types.CollectionData{types.CollectionDataFromCore(*data)})
	if err != nil {
		return fmt.Errorf("failed to create collection datas: %s", err)
	}

	col, err := m.source.Collection(height, data.Collection)
	if err != nil {
		return fmt.Errorf("failed to get collection from source: %s", err)
	}

	err = m.db.UpdateCollection(types.CollectionFromCore(col))
	if err != nil {
		return fmt.Errorf("failed to update collection: %s", err)
	}

	return nil
}

func (m *Module) RemoveCollectionDatas(height int64, datas []*tokenmanagertypes.CollectionDataIndex) (err error) {
	for _, data := range datas {
		err = m.db.RemoveCollectionData(tokenmanagertypes.CollectionDataKey(data))
		if err != nil {
			return fmt.Errorf("failed to remove collection data")
		}

		idx := types.NewCollectionDataIndex(data.Chain, data.Address)

		colData, err := m.source.CollectionData(height, *idx)
		if err != nil {
			return fmt.Errorf("failed to get collection data from source: %s", err)
		}

		err = m.updateCollection(colData.Collection)
		if err != nil {
			return fmt.Errorf("failed to update collection: %s", err)
		}
	}

	return nil
}

func (m *Module) updateCollection(index string) error {
	collection, err := m.source.Collection(0, index)
	if err != nil {
		return fmt.Errorf("failed to get collection from source: %s", err)
	}

	err = m.db.UpdateCollection(types.CollectionFromCore(collection))
	if err != nil {
		return fmt.Errorf("failed to update collection: %s", err)
	}

	return nil
}

func (m *Module) RemoveCollection(index string) (err error) {
	err = m.db.RemoveCollectionDataByCollection(index)
	if err != nil {
		return fmt.Errorf("failed to remove collection data: %s", err)
	}

	err = m.db.RemoveCollection(tokenmanagertypes.CollectionKey(index))
	if err != nil {
		return fmt.Errorf("failed to remove collection: %s", err)
	}

	return nil
}

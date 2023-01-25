package tokenmanager

import (
	"fmt"
	"gitlab.com/rarimo/bdjuno/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

func (m *Module) UpdateItems(items []*tokenmanagertypes.Item) (err error) {
	for _, item := range items {
		if item == nil {
			continue
		}

		err = m.db.UpdateItem(types.ItemFromCore(*item))
		if err != nil {
			return fmt.Errorf("failed to update item in tokenmanager: %s", err)
		}
	}

	return nil
}

func (m *Module) RemoveItems(indexes []*tokenmanagertypes.ItemIndex) (err error) {
	for _, index := range indexes {
		if index == nil {
			continue
		}

		err = m.db.RemoveItem(tokenmanagertypes.ItemKey(index))
		if err != nil {
			return fmt.Errorf("failed to remove item in tokenmanager: %s", err)
		}
	}

	return nil
}

func (m *Module) CreateCollection(
	index string,
	meta *tokenmanagertypes.CollectionMetadata,
	data []*tokenmanagertypes.CollectionData,
) error {
	coreCollection := tokenmanagertypes.Collection{
		Index: index,
		Meta:  meta,
		Data:  make([]*tokenmanagertypes.CollectionDataIndex, 0, len(data)),
	}

	for _, collectionData := range data {
		coreCollection.Data = append(coreCollection.Data, collectionData.Index)
	}

	collection := types.CollectionFromCore(coreCollection)

	return m.db.Transaction(func() error {
		err := m.db.SaveCollections([]types.Collection{collection})
		if err != nil {
			return fmt.Errorf("failed to create collection")
		}

		datas := make([]tokenmanagertypes.CollectionData, 0, len(data))
		for _, collectionData := range data {
			datas = append(datas, *collectionData)
		}

		err = m.saveCollectionDatas(datas)
		if err != nil {
			return fmt.Errorf("failed to create collection datas")
		}

		return nil
	})
}

func (m *Module) UpdateCollectionDatas(datas []*tokenmanagertypes.CollectionData) (err error) {
	for _, data := range datas {
		collectionData := types.CollectionDataFromCore(*data)
		err = m.db.UpdateCollectionData(collectionData)
		if err != nil {
			return fmt.Errorf("failed to update collection data")
		}
	}

	return nil
}

func (m *Module) CreateCollectionDatas(datas []*tokenmanagertypes.CollectionData) (err error) {
	list := make([]types.CollectionData, len(datas))

	for i, data := range datas {
		list[i] = types.CollectionDataFromCore(*data)
	}

	err = m.db.SaveCollectionDatas(list)
	if err != nil {
		return fmt.Errorf("failed to create collection datas")
	}

	return nil
}

func (m *Module) RemoveCollectionDatas(datas []*tokenmanagertypes.CollectionDataIndex) (err error) {
	for _, data := range datas {
		err = m.db.RemoveCollectionData(tokenmanagertypes.CollectionDataKey(data))
		if err != nil {
			return fmt.Errorf("failed to remove collection data")
		}
	}

	return nil
}

func (m *Module) RemoveCollection(index string) (err error) {
	err = m.db.RemoveCollection(tokenmanagertypes.CollectionKey(index))
	if err != nil {
		return fmt.Errorf("failed to remove collection")
	}

	return nil
}

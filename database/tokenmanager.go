package database

import (
	"encoding/json"
	"fmt"
	dbtypes "gitlab.com/rarimo/bdjuno/database/types"
	"gitlab.com/rarimo/bdjuno/types"
	"strings"
)

// SaveTokenManagerParams saves the given x/tokenmanager parameters inside the database
func (db *Db) SaveTokenManagerParams(params *types.TokenManagerParams) (err error) {
	paramsBz, err := json.Marshal(params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling tokenmanager params: %s", err)
	}

	stmt := `
INSERT INTO tokenmanager_params(params, height)
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE
	SET params = excluded.params,
		height = excluded.height
WHERE tokenmanager_params.height <= excluded.height
`
	_, err = db.Sql.Exec(
		stmt,
		string(paramsBz),
		params.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing tokenmanager params: %s", err)
	}

	return nil
}

func (db *Db) SaveCollections(collections []types.Collection) error {
	if len(collections) == 0 {
		return nil
	}

	collectionsQuery := `INSERT INTO collection (index, meta, data) VALUES `

	var collectionsParams []interface{}

	for i, collection := range collections {
		// Prepare the collection query
		vi := i * 3
		collectionsQuery += fmt.Sprintf("($%d, $%d, $%d),", vi+1, vi+2, vi+3)

		meta, err := json.Marshal(collection.Meta)
		if err != nil {
			return fmt.Errorf("error while marshaling meta: %s", err)
		}

		data, err := json.Marshal(collection.Data)
		if err != nil {
			return fmt.Errorf("error while marshaling data: %s", err)
		}

		collectionsParams = append(
			collectionsParams,
			collection.Index,
			string(meta),
			string(data),
		)
	}

	// Store the collections
	collectionsQuery = strings.TrimSuffix(collectionsQuery, ",") // Remove trailing ","
	collectionsQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(collectionsQuery, collectionsParams...)
	if err != nil {
		return fmt.Errorf("error while storing collections: %s", err)
	}

	return nil
}

func (db *Db) UpdateCollection(collection types.Collection) error {
	query := `UPDATE collection SET meta = $1, data = $2 WHERE index = $3`

	meta, err := json.Marshal(collection.Meta)
	if err != nil {
		return fmt.Errorf("error while marshaling meta: %s", err)
	}

	data, err := json.Marshal(collection.Data)
	if err != nil {
		return fmt.Errorf("error while marshaling data: %s", err)
	}

	index, err := json.Marshal(collection.Index)
	if err != nil {
		return fmt.Errorf("error while marshaling index: %s", err)
	}

	_, err = db.Sql.Exec(query, meta, data, index)
	if err != nil {
		return fmt.Errorf("error while updating collection: %s", err)
	}

	return nil
}

func (db *Db) SaveCollectionDatas(collectionDatas []types.CollectionData) error {
	if len(collectionDatas) == 0 {
		return nil
	}

	collectionDatasQuery := `INSERT INTO collection_data (index, index_key, collection, token_type, wrapped, decimals) VALUES `

	var collectionDatasParams []interface{}

	for i, collectionData := range collectionDatas {
		// Prepare the collection data query
		vi := i * 6
		collectionDatasQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6)

		index, err := json.Marshal(collectionData.Index)
		if err != nil {
			return fmt.Errorf("error while marshaling index: %s", err)
		}

		collectionDatasParams = append(
			collectionDatasParams,
			string(index),
			collectionData.IndexKey,
			collectionData.Collection,
			collectionData.TokenType,
			collectionData.Wrapped,
			collectionData.Decimals,
		)
	}

	// Store the collection datas
	collectionDatasQuery = strings.TrimSuffix(collectionDatasQuery, ",") // Remove trailing ","
	collectionDatasQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(collectionDatasQuery, collectionDatasParams...)
	if err != nil {
		return fmt.Errorf("error while storing collection datas: %s", err)
	}

	return nil
}

func (db *Db) UpdateCollectionData(data types.CollectionData) error {
	query := `UPDATE collection_data SET collection = $1, token_type = $2, wrapped = $3, decimals = $4 WHERE index_key = $5`

	_, err := db.Sql.Exec(query,
		data.Collection,
		data.TokenType,
		data.Wrapped,
		data.Decimals,
		data.IndexKey,
	)
	if err != nil {
		return fmt.Errorf("error while updating collection data: %s", err)
	}

	return nil
}

func (db *Db) SaveItems(items []types.Item) error {
	if len(items) == 0 {
		return nil
	}

	itemQuery := `INSERT INTO item (index, collection, meta, on_chain) VALUES `

	var itemParams []interface{}

	for i, item := range items {
		// Prepare the item data query
		vi := i * 4
		itemQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4)

		meta, err := json.Marshal(item.Meta)
		if err != nil {
			return fmt.Errorf("error while marshaling meta: %s", err)
		}

		onChain, err := json.Marshal(item.OnChain)
		if err != nil {
			return fmt.Errorf("error while marshaling on chain items: %s", err)
		}

		itemParams = append(
			itemParams,
			item.Index,
			item.Collection,
			string(meta),
			string(onChain),
		)
	}

	// Store the items
	itemQuery = strings.TrimSuffix(itemQuery, ",") // Remove trailing ","
	itemQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(itemQuery, itemParams...)
	if err != nil {
		return fmt.Errorf("error while storing items: %s", err)
	}

	return nil
}

func (db *Db) UpsertItem(item types.Item) error {
	meta, err := json.Marshal(item.Meta)
	if err != nil {
		return fmt.Errorf("error while marshaling meta: %s", err)
	}

	onChain, err := json.Marshal(item.OnChain)
	if err != nil {
		return fmt.Errorf("error while marshaling on chain items: %s", err)
	}

	stmt := `
INSERT INTO item (index, collection, meta, on_chain)
VALUES ($1, $2, $3, $4)
ON CONFLICT (index) DO UPDATE
	SET collection = excluded.collection, meta = excluded.meta, on_chain = excluded.on_chain
`

	_, err = db.Sql.Exec(
		stmt,
		item.Index,
		item.Collection,
		string(meta),
		string(onChain),
	)

	if err != nil {
		return fmt.Errorf("error while storing item: %s", err)
	}

	return nil
}

func (db *Db) RemoveItem(index string) error {
	stmt := `DELETE FROM item WHERE index = $1`
	_, err := db.Sql.Exec(stmt, index)
	if err != nil {
		return fmt.Errorf("error while deleting item: %s", err)
	}

	return nil
}

func (db *Db) GetItem(index string) (*types.Item, error) {
	stmt := `SELECT * FROM item WHERE index = $1`

	var items []dbtypes.ItemRow
	if err := db.Sqlx.Select(&items, stmt, index); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	row := items[0]

	var onChain []*types.OnChainItemIndex
	if err := json.Unmarshal([]byte(row.OnChain), &onChain); err != nil {
		return nil, fmt.Errorf("error while unmarshaling on chain items: %s", err)
	}

	var meta types.ItemMetadata
	if err := json.Unmarshal([]byte(row.Meta), &meta); err != nil {
		return nil, fmt.Errorf("error while unmarshaling meta: %s", err)
	}

	item := types.NewItem(
		row.Index,
		row.Collection,
		types.NewItemMetadata(
			meta.ImageUri,
			meta.ImageHash,
			meta.Seed,
			meta.Uri,
		),
		onChain,
	)

	return &item, nil
}

func (db *Db) RemoveCollectionData(indexKey []byte) error {
	stmt := `DELETE FROM collection_data WHERE index_key = $1`
	_, err := db.Sql.Exec(stmt, indexKey)
	if err != nil {
		return fmt.Errorf("error while deleting collection data: %s", err)
	}

	return nil
}

func (db *Db) RemoveCollectionDataByCollection(collection string) error {
	stmt := `DELETE FROM collection_data WHERE collection = $1`
	_, err := db.Sql.Exec(stmt, collection)
	if err != nil {
		return fmt.Errorf("error while deleting collection data: %s", err)
	}

	return nil
}

func (db *Db) RemoveCollection(indexKey []byte) error {
	stmt := `DELETE FROM collection WHERE index_key = $1`
	_, err := db.Sql.Exec(stmt, indexKey)
	if err != nil {
		return fmt.Errorf("error while deleting collection: %s", err)
	}

	return nil
}

func (db *Db) SaveOnChainItems(items []types.OnChainItem) error {
	if len(items) == 0 {
		return nil
	}

	query := `INSERT INTO on_chain_item (index, item) VALUES `

	var params []interface{}

	for i, item := range items {
		// Prepare the on chain item data query
		vi := i * 2
		query += fmt.Sprintf("($%d, $%d),", vi+1, vi+2)

		index, err := json.Marshal(item.Index)
		if err != nil {
			return fmt.Errorf("error while marshaling index: %s", err)
		}

		params = append(params, string(index), item.Item)
	}

	// Store the on chain items
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing on chain items: %s", err)
	}

	return nil
}

func (db *Db) RemoveOnChainItems(itemIndex string) error {
	stmt := `DELETE FROM on_chain_item WHERE item = $1`
	_, err := db.Sql.Exec(stmt, itemIndex)
	if err != nil {
		return fmt.Errorf("error while deleting on chain items: %s", err)
	}

	return nil
}

func (db *Db) SaveSeeds(seeds []types.Seed) error {
	if len(seeds) == 0 {
		return nil
	}

	query := `INSERT INTO seed (seed, item) VALUES `

	var params []interface{}

	for i, seed := range seeds {
		// Prepare the seed data query
		vi := i * 2
		query += fmt.Sprintf("($%d, $%d),", vi+1, vi+2)

		params = append(params, seed.Seed, seed.Item)
	}

	// Store the on chain seeds
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing on chain seeds: %s", err)
	}

	return nil
}

func (db *Db) UpsertSeed(seed types.Seed) error {
	stmt := `INSERT INTO seed (seed, item) VALUES ($1, $2) ON CONFLICT (item) DO UPDATE SET seed = excluded.seed`

	_, err := db.Sql.Exec(stmt, seed.Seed, seed.Item)

	if err != nil {
		return fmt.Errorf("error while storing item: %s", err)
	}

	return nil
}

func (db *Db) RemoveSeed(seed string) error {
	stmt := `DELETE FROM seed WHERE seed = $1`
	_, err := db.Sql.Exec(stmt, seed)
	if err != nil {
		return fmt.Errorf("error while deleting seed: %s", err)
	}

	return nil
}

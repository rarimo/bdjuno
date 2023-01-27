package database

import (
	"fmt"
	"github.com/lib/pq"
	dbtypes "gitlab.com/rarimo/bdjuno/database/types"
	"gitlab.com/rarimo/bdjuno/types"
)

// SaveTokenManagerParams saves the given x/tokenmanager parameters inside the database
func (db *Db) SaveTokenManagerParams(params *types.TokenManagerParams) (err error) {
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
		params.Params,
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

	collectionsQuery := `INSERT INTO collection (index, item_key, meta, data) VALUES `

	var collectionsParams []interface{}

	for i, collection := range collections {
		// Prepare the collection query
		vi := i * 4
		collectionsQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4)

		collectionsParams = append(
			collectionsParams,
			collection.Index,
			collection.IndexKey,
			collection.Meta,
			pq.Array(collection.Data),
		)
	}

	// Store the collections
	collectionsQuery = collectionsQuery[:len(collectionsQuery)-1] // Remove trailing ","
	collectionsQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(collectionsQuery, collectionsParams...)
	if err != nil {
		return fmt.Errorf("error while storing collections: %s", err)
	}

	return nil
}

func (db *Db) UpdateCollection(collection types.Collection) error {
	query := `UPDATE collection SET meta = $1, data = $2 WHERE index = $3`

	_, err := db.Sql.Exec(query, collection.Meta, pq.Array(collection.Data), collection.Index)
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

		collectionDatasParams = append(
			collectionDatasParams,
			collectionData.Index,
			collectionData.IndexKey,
			collectionData.Collection,
			collectionData.TokenType,
			collectionData.Wrapped,
			collectionData.Decimals,
		)
	}

	// Store the collection datas
	collectionDatasQuery = collectionDatasQuery[:len(collectionDatasQuery)-1] // Remove trailing ","
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

		itemParams = append(
			itemParams,
			item.Index,
			item.Collection,
			item.Meta,
			pq.Array(item.OnChain),
		)
	}

	// Store the items
	itemQuery = itemQuery[:len(itemQuery)-1] // Remove trailing ","
	itemQuery += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(itemQuery, itemParams...)
	if err != nil {
		return fmt.Errorf("error while storing items: %s", err)
	}

	return nil
}

func (db *Db) UpdateItem(item types.Item) error {
	query := `UPDATE item SET meta = $1, on_chain = $2 WHERE index = $3`

	_, err := db.Sql.Exec(query, item.Meta, pq.Array(item.OnChain), item.Index)
	if err != nil {
		return fmt.Errorf("error while updating item: %s", err)
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
	meta := row.Meta

	onChain := make([]*types.OnChainItemIndex, len(row.OnChain))
	for i, onChainItem := range row.OnChain {
		onChain[i] = types.NewOnChainItemIndex(onChainItem.Chain, onChainItem.Address, onChainItem.TokenID)
	}

	item := types.NewItem(
		row.Index,
		row.Collection,
		types.NewItemMetadata(
			meta.ImageUri,
			meta.ImageHash,
			meta.Seed,
			meta.Name,
			meta.Symbol,
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

		params = append(params, item.Index, item.Item)
	}

	// Store the on chain items
	query = query[:len(query)-1] // Remove trailing ","
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
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing on chain seeds: %s", err)
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

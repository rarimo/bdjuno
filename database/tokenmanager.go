package database

import (
	"fmt"
	"github.com/lib/pq"
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
		pq.Array(params.Params),
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

func (db *Db) SaveCollectionDatas(collectionDatas []types.CollectionData) error {
	if len(collectionDatas) == 0 {
		return nil
	}

	collectionDatasQuery := `INSERT INTO collection_data (index, index_key, collection, token_type, wrapped, decimals) VALUES `

	var collectionDatasParams []interface{}

	for i, collectionData := range collectionDatas {
		// Prepare the collection data query
		vi := i * 5
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

	itemQuery := `INSERT INTO item (index, item_key, meta, chain_params) VALUES `

	var itemParams []interface{}

	for i, item := range items {
		// Prepare the item data query
		vi := i * 3
		itemQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4)

		itemParams = append(
			itemParams,
			item.Index,
			item.IndexKey,
			item.Meta,
			item.ChainParams,
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
	query := `UPDATE item SET meta = $1, chain_params = $2 WHERE index_key = $3`

	_, err := db.Sql.Exec(query,
		item.Meta,
		pq.Array(item.ChainParams),
		item.IndexKey,
	)
	if err != nil {
		return fmt.Errorf("error while updating item: %s", err)
	}

	return nil
}

func (db *Db) RemoveItem(indexKey []byte) error {
	stmt := `DELETE FROM item WHERE index_key = $1`
	_, err := db.Sql.Exec(stmt, indexKey)
	if err != nil {
		return fmt.Errorf("error while deleting item: %s", err)
	}

	return nil
}

func (db *Db) RemoveCollectionData(indexKey []byte) error {
	stmt := `DELETE FROM collection_data WHERE index_key = $1`
	_, err := db.Sql.Exec(stmt, indexKey)
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

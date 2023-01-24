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

	collectionsQuery := `INSERT INTO collection (index, name, symbol, metadata_uri, data) VALUES `

	var collectionsParams []interface{}

	for i, collection := range collections {
		// Prepare the collection query
		vi := i * 5
		collectionsQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5)

		collectionsParams = append(
			collectionsParams,
			collection.Index,
			collection.Name,
			collection.Symbol,
			collection.MetadataURI,
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

	collectionDatasQuery := `INSERT INTO collection_data (index, collection, token_type, wrapped, decimals) VALUES `

	var collectionDatasParams []interface{}

	for i, collectionData := range collectionDatas {
		// Prepare the collection data query
		vi := i * 5
		collectionDatasQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5)

		collectionDatasParams = append(
			collectionDatasParams,
			collectionData.Index,
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

func (db *Db) SaveItems(items []types.Item) error {
	if len(items) == 0 {
		return nil
	}

	itemQuery := `INSERT INTO item (index, meta, chain_params) VALUES `

	var itemParams []interface{}

	for i, item := range items {
		// Prepare the item data query
		vi := i * 3
		itemQuery += fmt.Sprintf("($%d, $%d, $%d),", vi+1, vi+2, vi+3)

		itemParams = append(
			itemParams,
			item.Index,
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

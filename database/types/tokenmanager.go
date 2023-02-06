package types

// TokenManagerParamsRow represents a single row of the "tokenmanager_params" table
type TokenManagerParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// CollectionMetadata represents the information stored inside the database about a collection metadata
type CollectionMetadata struct {
	Name        string `db:"name"`
	Symbol      string `db:"symbol"`
	MetadataURI string `db:"metadata_uri"`
}

// CollectionDataRow represents a single row of the "collection_data" table
type CollectionDataRow struct {
	Index      *string `db:"index"`
	IndexKey   string  `db:"index_key"`
	Collection string  `db:"collection"`
	TokenType  int32   `db:"tokenType"`
	Wrapped    bool    `db:"wrapped"`
	Decimals   uint32  `db:"decimals"`
}

// CollectionRow represents a single row of the "collection" table
type CollectionRow struct {
	Index string `db:"index"`
	Meta  string `db:"meta"`
	Data  string `db:"data"`
}

// OnChainItemIndex represents the information stored inside the database about an item index
type OnChainItemIndex struct {
	Chain   string `db:"chain"`
	Address string `db:"address"`
	TokenID string `db:"token_id"`
}

// ItemMetadata represents the information stored inside the database about an item metadata
type ItemMetadata struct {
	ImageUri  string `db:"image_uri"`
	ImageHash string `db:"image_hash"`
	Seed      string `db:"seed"`
	Uri       string `db:"uri"`
}

// ItemRow represents a single row of the "item" table
type ItemRow struct {
	Index      string `db:"index"`
	Collection string `db:"collection"`
	Meta       string `db:"meta"`
	OnChain    string `db:"on_chain"`
}

// OnChainItemRow represents a single row of the "on_chain_item" table
type OnChainItemRow struct {
	Index string `db:"index"`
	Item  string `db:"item"`
}

// Seed represents a single row of the "seed" table
type Seed struct {
	Seed string `db:"seed"`
	Item string `db:"item"`
}

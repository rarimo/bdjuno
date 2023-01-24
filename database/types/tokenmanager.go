package types

// NetworkTypeBinding represents the information stored inside the database about a network type binding
type NetworkTypeBinding struct {
	CoreType  uint32 `db:"core_type"`
	SaverType uint32 `db:"saver_type"`
}

// NetworkParams represents the information stored inside the database about a network params
type NetworkParams struct {
	Name     string                `db:"name"`
	Contract string                `db:"contract"`
	Types    []*NetworkTypeBinding `db:"types"`
	Type     int32                 `db:"type"`
}

type TokenManagerParamsInner struct {
	Networks []NetworkParams `db:"networks"`
}

// TokenManagerParamsRow represents a single row of the "tokenmanager_params" table
type TokenManagerParamsRow struct {
	OneRowID bool                    `db:"one_row_id"`
	Params   TokenManagerParamsInner `db:"params"`
	Height   int64                   `db:"height"`
}

// CollectionDataIndex represents the information stored inside the database about a collection data index
type CollectionDataIndex struct {
	Chain   string `db:"chain"`
	Address string `db:"address"`
}

// CollectionMetadata represents the information stored inside the database about a collection metadata
type CollectionMetadata struct {
	Name        string `db:"name"`
	Symbol      string `db:"symbol"`
	MetadataURI string `db:"metadata_uri"`
}

// CollectionDataRow represents a single row of the "collection_data" table
type CollectionDataRow struct {
	Index      *CollectionDataIndex `db:"index"`
	Collection string               `db:"collection"`
	TokenType  int32                `db:"tokenType"`
	Wrapped    bool                 `db:"wrapped"`
	Decimals   uint32               `db:"decimals"`
}

// CollectionRow represents a single row of the "collection" table
type CollectionRow struct {
	Index string                 `db:"index"`
	Meta  *CollectionMetadata    `db:"meta"`
	Data  []*CollectionDataIndex `db:"data"`
}

// ItemIndex represents the information stored inside the database about an item index
type ItemIndex struct {
	Collection string `db:"collection"`
	Name       string `db:"name"`
	Symbol     string `db:"symbol"`
	Uri        string `db:"uri"`
}

// ItemMetadata represents the information stored inside the database about an item metadata
type ItemMetadata struct {
	ImageUri  string `db:"image_uri"`
	ImageHash string `db:"image_hash"`
	Seed      string `db:"seed"`
}

// ItemChainParams represents the information stored inside the database about an item chain params
type ItemChainParams struct {
	Chain   string `db:"chain"`
	TokenID string `db:"tokenID"`
}

// ItemRow represents a single row of the "item" table
type ItemRow struct {
	Index       *ItemIndex         `db:"index"`
	Meta        *ItemMetadata      `db:"meta"`
	ChainParams []*ItemChainParams `db:"chain_params"`
}

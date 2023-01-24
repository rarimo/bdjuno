package types

import (
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

// NetworkTypeBinding contains the data of the x/tokenmanager network core and saver type bindings
type NetworkTypeBinding struct {
	CoreType  uint32 `json:"core_type,omitempty" yaml:"core_type,omitempty"`
	SaverType uint32 `json:"saver_type,omitempty" yaml:"saver_type,omitempty"`
}

// NetworkParams contains the data of the x/tokenmanager network params
type NetworkParams struct {
	Name     string                        `json:"name,omitempty" yaml:"name,omitempty"`
	Contract string                        `json:"contract,omitempty" yaml:"contract,omitempty"`
	Types    []*NetworkTypeBinding         `json:"types,omitempty" yaml:"types,omitempty"`
	Type     tokenmanagertypes.NetworkType `json:"type,omitempty" yaml:"type,omitempty"`
}

type TokenManagerParamsInner struct {
	Networks []NetworkParams `json:"networks,omitempty" yaml:"networks,omitempty"`
}

// TokenManagerParams contains the data of the x/tokenmanager module params instance
type TokenManagerParams struct {
	Params TokenManagerParamsInner `json:"params,omitempty" yaml:"params,omitempty"`
	Height int64                   `json:"height,omitempty" yaml:"height,omitempty"`
}

// NewTokenManagerParams allows to build a new TokenManagerParams instance
func NewTokenManagerParams(params tokenmanagertypes.Params, height int64) *TokenManagerParams {
	networks := make([]NetworkParams, 0)

	for _, network := range params.Networks {
		if network == nil {
			continue
		}

		networkTypes := make([]*NetworkTypeBinding, len(network.Types))

		for _, networkType := range network.Types {
			if networkType == nil {
				continue
			}

			networkTypes = append(networkTypes, &NetworkTypeBinding{
				CoreType:  networkType.CoreType,
				SaverType: networkType.SaverType,
			})
		}

		networks = append(networks, NetworkParams{
			Name:     network.Name,
			Contract: network.Contract,
			Type:     network.Type,
			Types:    networkTypes,
		})

	}

	return &TokenManagerParams{
		Params: TokenManagerParamsInner{Networks: networks},
		Height: height,
	}
}

// CollectionDataIndex contains the data of the x/tokenmanager collection data index
type CollectionDataIndex struct {
	Chain   string `json:"chain,omitempty" yaml:"chain,omitempty"`
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
}

// NewCollectionDataIndex allows to build a new CollectionDataIndex instance
func NewCollectionDataIndex(chain, address string) *CollectionDataIndex {
	return &CollectionDataIndex{
		Chain:   chain,
		Address: address,
	}
}

// CollectionData contains the data of the x/tokenmanager collection data instance
type CollectionData struct {
	Index      *CollectionDataIndex   `json:"index,omitempty" yaml:"index,omitempty"`
	Collection string                 `json:"collection,omitempty" yaml:"collection,omitempty"`
	TokenType  tokenmanagertypes.Type `json:"token_type,omitempty" yaml:"tokenType,omitempty"`
	Wrapped    bool                   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
	Decimals   uint32                 `json:"decimals,omitempty" yaml:"decimals,omitempty"`
}

// NewCollectionData allows to build a new CollectionData instance
func NewCollectionData(
	index *CollectionDataIndex,
	collection string,
	tokenType tokenmanagertypes.Type,
	wrapped bool,
	decimals uint32,
) CollectionData {
	return CollectionData{
		Index:      index,
		Collection: collection,
		TokenType:  tokenType,
		Wrapped:    wrapped,
		Decimals:   decimals,
	}
}

// CollectionMetadata contains the data of the x/tokenmanager collection metadata related to collection instance
type CollectionMetadata struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Symbol      string `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	MetadataURI string `json:"metadata_uri,omitempty" yaml:"metadata_uri,omitempty"`
}

// NewCollectionMetadata allows to build a new CollectionMetadata instance
func NewCollectionMetadata(name, symbol, metadataUri string) *CollectionMetadata {
	return &CollectionMetadata{
		Name:        name,
		Symbol:      symbol,
		MetadataURI: metadataUri,
	}
}

// Collection contains the data of the x/tokenmanager collection instance
type Collection struct {
	Index string                 `json:"index,omitempty" yaml:"index,omitempty"`
	Meta  *CollectionMetadata    `json:"meta,omitempty" yaml:"data,omitempty"`
	Data  []*CollectionDataIndex `json:"data,omitempty" yaml:"data,omitempty"`
}

// NewCollection allows to build a new Collection instance
func NewCollection(index string, meta *CollectionMetadata, data []*CollectionDataIndex) Collection {
	return Collection{
		Index: index,
		Meta:  meta,
		Data:  data,
	}
}

// ItemIndex contains the data of the x/tokenmanager item index
type ItemIndex struct {
	Collection string `json:"collection,omitempty" yaml:"collection,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Symbol     string `json:"symbol,omitempty" yaml:"symbol,omitempty"`
	Uri        string `json:"uri,omitempty" yaml:"uri,omitempty"`
}

// NewItemIndex allows to build a new ItemIndex instance
func NewItemIndex(collection, name, symbol, uri string) *ItemIndex {
	return &ItemIndex{
		Collection: collection,
		Name:       name,
		Symbol:     symbol,
		Uri:        uri,
	}
}

// ItemMetadata contains the data of the x/tokenmanager item metadata
type ItemMetadata struct {
	ImageUri  string `json:"image_uri,omitempty" yaml:"image_uri,omitempty"`
	ImageHash string `json:"image_hash,omitempty" yaml:"image_hash,omitempty"`
	Seed      string `json:"seed,omitempty" yaml:"seed,omitempty"`
}

// NewItemMetadata allows to build a new ItemMetadata instance
func NewItemMetadata(imageUri, imageHash, seed string) *ItemMetadata {
	return &ItemMetadata{
		ImageUri:  imageUri,
		ImageHash: imageHash,
		Seed:      seed,
	}
}

// ItemChainParams contains the data of the x/tokenmanager item chain params
type ItemChainParams struct {
	Chain   string `json:"chain,omitempty" yaml:"chain,omitempty"`
	TokenID string `json:"tokenID,omitempty" yaml:"tokenID,omitempty"`
}

// NewItemChainParams allows to build a new ItemChainParams instance
func NewItemChainParams(chain, tokenID string) *ItemChainParams {
	return &ItemChainParams{
		Chain:   chain,
		TokenID: tokenID,
	}
}

// Item contains the data of the x/tokenmanager item instance
type Item struct {
	Index       *ItemIndex         `json:"index,omitempty" yaml:"index,omitempty"`
	Meta        *ItemMetadata      `json:"meta,omitempty" yaml:"meta,omitempty"`
	ChainParams []*ItemChainParams `json:"chain_params,omitempty" yaml:"chain_params,omitempty"`
}

// NewItem allows to build a new Item instance
func NewItem(index *ItemIndex, meta *ItemMetadata, chainParams []*ItemChainParams) Item {
	return Item{
		Index:       index,
		Meta:        meta,
		ChainParams: chainParams,
	}
}

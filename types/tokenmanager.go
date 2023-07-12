package types

import (
	"encoding/json"
	"fmt"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

//--------------------------------------------------------

// Network contains the data of the x/tokenmanager network params
type Network struct {
	Name   string                        `json:"name,omitempty" yaml:"name,omitempty"`
	Type   tokenmanagertypes.NetworkType `json:"type" yaml:"type"`
	Params []NetworkParams               `json:"params,omitempty" yaml:"params,omitempty"`
}

type NetworkParams struct {
	Type    tokenmanagertypes.NetworkParamType `json:"name,omitempty" yaml:"name,omitempty"`
	Details json.RawMessage                    `json:"details,omitempty" yaml:"details,omitempty"`
}

type TokenManagerParamsInner struct {
	Networks []Network `json:"networks,omitempty" yaml:"networks,omitempty"`
}

// TokenManagerParams contains the data of the x/tokenmanager module params instance
type TokenManagerParams struct {
	Params TokenManagerParamsInner `json:"params,omitempty" yaml:"params,omitempty"`
	Height int64                   `json:"height,omitempty" yaml:"height,omitempty"`
}

// NewTokenManagerParams allows to build a new TokenManagerParams instance
func NewTokenManagerParams(params tokenmanagertypes.Params, height int64) (*TokenManagerParams, error) {
	networks := make([]Network, 0)

	for _, network := range params.Networks {
		n := Network{
			Name:   network.Name,
			Type:   network.Type,
			Params: make([]NetworkParams, len(network.Params)),
		}

		for i, networkParams := range network.Params {
			details, err := json.Marshal(networkParams.Details)
			if err != nil {
				return nil, fmt.Errorf("error while marshalling network params details: %s", err)
			}
			n.Params[i] = NetworkParams{
				Type:    networkParams.Type,
				Details: details,
			}
		}

		networks = append(networks, n)

	}

	return &TokenManagerParams{
		Params: TokenManagerParamsInner{
			Networks: networks,
		},
		Height: height,
	}, nil
}

//--------------------------------------------------------

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
	IndexKey   string                 `json:"index_key" yaml:"index_key"`
	Collection string                 `json:"collection,omitempty" yaml:"collection,omitempty"`
	TokenType  tokenmanagertypes.Type `json:"token_type" yaml:"tokenType"`
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
	indexKey := tokenmanagertypes.CollectionDataKey(&tokenmanagertypes.CollectionDataIndex{
		Chain:   index.Chain,
		Address: index.Address,
	})

	return CollectionData{
		Index:      index,
		IndexKey:   string(indexKey),
		Collection: collection,
		TokenType:  tokenType,
		Wrapped:    wrapped,
		Decimals:   decimals,
	}
}

// CollectionDataFromCore allows to build a new CollectionData instance from tokenmanager.CollectionData instance
func CollectionDataFromCore(data tokenmanagertypes.CollectionData) CollectionData {
	return NewCollectionData(
		NewCollectionDataIndex(data.Index.Chain, data.Index.Address),
		data.Collection,
		data.TokenType,
		data.Wrapped,
		data.Decimals,
	)
}

//--------------------------------------------------------

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

//--------------------------------------------------------

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

// CollectionFromCore allows to build a new Collection instance from tokenmanager.Collection instance
func CollectionFromCore(collection tokenmanagertypes.Collection) Collection {
	indexes := make([]*CollectionDataIndex, 0)

	for _, data := range collection.Data {
		if data == nil {
			continue
		}

		indexes = append(indexes, NewCollectionDataIndex(data.Chain, data.Address))
	}

	return NewCollection(
		collection.Index,
		NewCollectionMetadata(
			collection.Meta.Name,
			collection.Meta.Symbol,
			collection.Meta.MetadataURI,
		),
		indexes,
	)
}

//--------------------------------------------------------

// OnChainItemIndex contains the data of the x/tokenmanager on chain item index
type OnChainItemIndex struct {
	Chain   string `json:"chain,omitempty" yaml:"chain,omitempty"`
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
	TokenID string `json:"token_id,omitempty" yaml:"token_id,omitempty"`
}

// OnChainItemIndexFromCore allows to build a new OnChainItemIndex instance from tokenmanager.OnChainItemIndex instance
func OnChainItemIndexFromCore(index *tokenmanagertypes.OnChainItemIndex) *OnChainItemIndex {
	return &OnChainItemIndex{
		Chain:   index.Chain,
		Address: index.Address,
		TokenID: index.TokenID,
	}
}

//--------------------------------------------------------

// ItemMetadata contains the data of the x/tokenmanager item metadata
type ItemMetadata struct {
	ImageUri  string `json:"image_uri,omitempty" yaml:"image_uri,omitempty"`
	ImageHash string `json:"image_hash,omitempty" yaml:"image_hash,omitempty"`
	Seed      string `json:"seed,omitempty" yaml:"seed,omitempty"`
	Uri       string `json:"uri,omitempty" yaml:"uri,omitempty"`
}

// NewItemMetadata allows to build a new ItemMetadata instance
func NewItemMetadata(imageUri, imageHash, seed, uri string) *ItemMetadata {
	return &ItemMetadata{
		ImageUri:  imageUri,
		ImageHash: imageHash,
		Seed:      seed,
		Uri:       uri,
	}
}

// ItemMetadataFromCore allows to build a new ItemMetadata instance from tokenmanager.ItemMetadata instance
func ItemMetadataFromCore(meta *tokenmanagertypes.ItemMetadata) *ItemMetadata {
	if meta != nil {
		return &ItemMetadata{
			ImageUri:  meta.ImageUri,
			ImageHash: meta.ImageHash,
			Seed:      meta.Seed,
			Uri:       meta.Uri,
		}
	}

	return nil
}

//--------------------------------------------------------

// Item contains the data of the x/tokenmanager item instance
type Item struct {
	Index      string              `json:"index,omitempty" yaml:"index,omitempty"`
	Collection string              `json:"collection,omitempty" yaml:"collection,omitempty"`
	Meta       *ItemMetadata       `json:"meta,omitempty" yaml:"meta,omitempty"`
	OnChain    []*OnChainItemIndex `json:"on_chain,omitempty" yaml:"on_chain,omitempty"`
}

// NewItem allows to build a new Item instance
func NewItem(index, collection string, meta *ItemMetadata, onChain []*OnChainItemIndex) Item {
	return Item{
		Index:      index,
		Collection: collection,
		Meta:       meta,
		OnChain:    onChain,
	}
}

// ItemFromCore allows to build a new Item instance from tokenmanager.Item instance
func ItemFromCore(item tokenmanagertypes.Item) Item {
	indexes := make([]*OnChainItemIndex, 0)

	for _, onChain := range item.OnChain {
		if onChain == nil {
			continue
		}

		indexes = append(indexes, OnChainItemIndexFromCore(onChain))
	}

	meta := ItemMetadataFromCore(&item.Meta)

	return NewItem(
		item.Index,
		item.Collection,
		meta,
		indexes,
	)
}

//--------------------------------------------------------

// OnChainItem contains the data of the x/tokenmanager on chain item instance
type OnChainItem struct {
	Index *OnChainItemIndex `json:"index,omitempty" yaml:"index,omitempty"`
	Item  string            `json:"item,omitempty" yaml:"item,omitempty"`
}

// OnChainItemFromCore allows to build a new OnChainItem instance from tokenmanager.OnChainItem instance
func OnChainItemFromCore(item tokenmanagertypes.OnChainItem) OnChainItem {
	return OnChainItem{
		Index: OnChainItemIndexFromCore(item.Index),
		Item:  item.Item,
	}
}

//--------------------------------------------------------

// Seed contains the data of the x/tokenmanager seed instance
type Seed struct {
	Seed string `json:"seed,omitempty" yaml:"seed,omitempty"`
	Item string `json:"item,omitempty" yaml:"item,omitempty"`
}

// NewSeed allows to build a new Seed instance
func NewSeed(seed, item string) Seed {
	return Seed{
		Seed: seed,
		Item: item,
	}
}

// SeedFromCore allows to build a new Seed instance from tokenmanager.Seed instance
func SeedFromCore(seed tokenmanagertypes.Seed) Seed {
	return Seed{
		Seed: seed.Seed,
		Item: seed.Item,
	}
}

//--------------------------------------------------------

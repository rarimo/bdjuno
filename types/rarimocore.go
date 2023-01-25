package types

import (
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

// Party contains the data of the x/rarimocore module signer instance
type Party struct {
	Account  string `json:"account,omitempty" yaml:"account,omitempty"`
	PubKey   string `json:"pub_key,omitempty" yaml:"pub_key,omitempty"`
	Address  string `json:"address,omitempty" yaml:"address,omitempty"`
	Verified bool   `json:"verified,omitempty" yaml:"verified,omitempty"`
}

// NewParty allows to build a new Party
func NewParty(p rarimocoretypes.Party) Party {
	return Party{
		Account:  p.Account,
		PubKey:   p.PubKey,
		Address:  p.Address,
		Verified: p.Verified,
	}
}

// RarimoCoreParams contains the data of the x/rarimocore module params instance
type RarimoCoreParams struct {
	KeyECDSA                  string   `json:"key_ecdsa,omitempty" yaml:"key_ecdsa,omitempty"`
	Threshold                 uint64   `json:"threshold,omitempty" yaml:"threshold,omitempty"`
	IsUpdateRequired          bool     `json:"is_update_required,omitempty" yaml:"is_update_required,omitempty"`
	LastSignature             string   `json:"last_signature,omitempty" yaml:"last_signature,omitempty"`
	Parties                   []string `json:"parties,omitempty" yaml:"parties,omitempty"`
	Height                    int64    `json:"height,omitempty" yaml:"height,omitempty"`
	AvailableResignBlockDelta uint64   `json:"available_resign_block_delta,omitempty" yaml:"available_resign_block_delta,omitempty"`
}

// NewRarimoCoreParams allows to build a new RarimoCoreParams instance
func NewRarimoCoreParams(p rarimocoretypes.Params, height int64) *RarimoCoreParams {
	parties := make([]string, len(p.Parties))
	for i, party := range p.Parties {
		parties[i] = party.Account
	}
	return &RarimoCoreParams{
		KeyECDSA:                  p.KeyECDSA,
		Threshold:                 p.Threshold,
		IsUpdateRequired:          p.IsUpdateRequired,
		LastSignature:             p.LastSignature,
		AvailableResignBlockDelta: p.AvailableResignBlockDelta,
		Parties:                   parties,
		Height:                    height,
	}
}

// Operation represents a single operation instance
type Operation struct {
	Index         string                 `json:"index,omitempty" yaml:"index,omitempty"`
	OperationType rarimocoretypes.OpType `json:"operation_type,omitempty" yaml:"operation_type,omitempty"`
	Signed        bool                   `json:"signed,omitempty" yaml:"signed,omitempty"`
	Approved      bool                   `json:"approved,omitempty" yaml:"approved,omitempty"`
	Creator       string                 `json:"creator,omitempty" yaml:"creator,omitempty"`
	Timestamp     int64                  `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
}

// NewOperation allows to build a new Operation instance
func NewOperation(index string, opType int32, signed, approved bool, creator string, timestamp uint64) Operation {
	return Operation{
		Index:         index,
		OperationType: rarimocoretypes.OpType(opType),
		Approved:      approved,
		Signed:        signed,
		Creator:       creator,
		Timestamp:     int64(timestamp),
	}
}

// OperationFromCore allows to build a new Operation instance from a rarimocoretypes.Operation instance
func OperationFromCore(operation rarimocoretypes.Operation) Operation {
	return Operation{
		Index:         operation.Index,
		OperationType: operation.OperationType,
		Approved:      operation.Approved,
		Signed:        operation.Signed,
		Creator:       operation.Creator,
		Timestamp:     int64(operation.Timestamp),
	}
}

// Transfer represents a single transfer instance
type Transfer struct {
	OperationIndex string           `json:"operation_index,omitempty" yaml:"operation_index,omitempty"`
	Origin         string           `json:"origin,omitempty" yaml:"origin,omitempty"`
	Tx             string           `json:"tx,omitempty" yaml:"tx,omitempty"`
	EventID        string           `json:"event_id,omitempty" yaml:"event_id,omitempty"`
	Receiver       string           `json:"receiver,omitempty" yaml:"receiver,omitempty"`
	Amount         string           `json:"amount,omitempty" yaml:"amount,omitempty"`
	BundleData     string           `json:"bundle_data,omitempty" yaml:"bundle_data,omitempty"`
	BundleSalt     string           `json:"bundle_salt,omitempty" yaml:"bundle_salt,omitempty"`
	ItemIndexKey   []byte           `json:"item_index_key,omitempty" yaml:"item_index_key,omitempty"`
	FromChain      *ItemChainParams `json:"from_chain,omitempty" yaml:"from_chain,omitempty"`
	ToChain        *ItemChainParams `json:"to_chain,omitempty" yaml:"to_chain,omitempty"`
	ItemIndex      *ItemIndex       `json:"item_index,omitempty" yaml:"item_index,omitempty"`
	ItemMeta       *ItemMetadata    `json:"item_meta,omitempty" yaml:"item_meta,omitempty"`
}

// NewTransfer allows to build a new Transfer instance
func NewTransfer(operationIndex string, t rarimocoretypes.Transfer) Transfer {
	return Transfer{
		OperationIndex: operationIndex,
		Origin:         t.Origin,
		Tx:             t.Tx,
		EventID:        t.EventId,
		Receiver:       t.Receiver,
		Amount:         t.Amount,
		BundleData:     t.BundleData,
		BundleSalt:     t.BundleSalt,
		ItemIndexKey:   tokenmanagertypes.ItemKey(t.Item),
		ItemIndex:      ItemIndexFromCore(t.Item),
		ItemMeta:       ItemMetadataFromCore(t.Meta),
		FromChain:      ItemChainParamsFromCore(t.From),
		ToChain:        ItemChainParamsFromCore(t.To),
	}
}

// ChangeParties represents a single change parties instance
type ChangeParties struct {
	OperationIndex string   `json:"operation_index,omitempty" yaml:"operation_index,omitempty"`
	Parties        []string `json:"parties,omitempty" yaml:"parties,omitempty"`
	NewPublicKey   string   `json:"new_public_key,omitempty" yaml:"new_public_key,omitempty"`
	Signature      string   `json:"signature,omitempty" yaml:"signature,omitempty"`
}

// NewChangeParties allows to build a new ChangeParties instance
func NewChangeParties(operationIndex string, c rarimocoretypes.ChangeParties) ChangeParties {
	parties := make([]string, len(c.Parties))
	for i := range c.Parties {
		parties[i] = c.Parties[i].Account
	}

	return ChangeParties{
		OperationIndex: operationIndex,
		Parties:        parties,
		NewPublicKey:   c.NewPublicKey,
		Signature:      c.Signature,
	}
}

// Confirmation represents a single confirmation instance
type Confirmation struct {
	Root           string   `json:"root,omitempty" yaml:"root,omitempty"`
	Indexes        []string `json:"indexes,omitempty" yaml:"indexes,omitempty"`
	SignatureECDSA string   `json:"signatureECDSA,omitempty" yaml:"signatureECDSA,omitempty"`
	Creator        string   `json:"creator,omitempty" yaml:"creator,omitempty"`
}

// NewConfirmation allows to build a new Confirmation instance
func NewConfirmation(c rarimocoretypes.Confirmation) Confirmation {
	return Confirmation{
		Root:           c.Root,
		Indexes:        c.Indexes,
		SignatureECDSA: c.SignatureECDSA,
		Creator:        c.Creator,
	}
}

// RarimoCoreVote represents a single vote instance
type RarimoCoreVote struct {
	Operation string                   `json:"operation,omitempty" yaml:"operation,omitempty"`
	Validator string                   `json:"validator,omitempty" yaml:"validator,omitempty"`
	Vote      rarimocoretypes.VoteType `json:"vote,omitempty" yaml:"vote,omitempty"`
}

// NewRarimoCoreVote allows to build a new RarimoCoreVote instance
func NewRarimoCoreVote(operation, validator string, vote rarimocoretypes.VoteType) RarimoCoreVote {
	return RarimoCoreVote{
		Operation: operation,
		Validator: validator,
		Vote:      vote,
	}
}

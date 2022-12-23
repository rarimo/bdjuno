package types

import (
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

// PartyRow represents a single row of the "parties" table
type PartyRow struct {
	Account  string `db:"account"`
	PubKey   string `db:"pub_key"`
	Address  string `db:"address"`
	Verified bool   `db:"verified"`
}

func NewPartyRow(
	account string,
	pubKey string,
	address string,
	verified bool,
) PartyRow {
	return PartyRow{
		Account:  account,
		PubKey:   pubKey,
		Address:  address,
		Verified: verified,
	}
}

// RarimoCoreParamsRow represents a single row of the "rarimocore_params" table
type RarimoCoreParamsRow struct {
	OneRowID         bool     `db:"one_row_id"`
	KeyECDSA         string   `db:"key_ecdsa"`
	Threshold        uint64   `db:"threshold"`
	IsUpdateRequired bool     `db:"is_update_required"`
	LastSignature    string   `db:"last_signature"`
	Parties          []string `db:"parties"`
	Height           int64    `db:"height"`
}

type OperationRow struct {
	Index         string                 `db:"index"`
	OperationType rarimocoretypes.OpType `db:"operation_type"`
	Signed        bool                   `db:"signed"`
	Creator       string                 `db:"creator"`
	Timestamp     int64                  `db:"timestamp"`
}

func NewOperationRow(
	index string,
	operationType rarimocoretypes.OpType,
	signed bool,
	creator string,
	timestamp int64,
) OperationRow {
	return OperationRow{
		Index:         index,
		OperationType: operationType,
		Signed:        signed,
		Creator:       creator,
		Timestamp:     timestamp,
	}
}

type TransferRow struct {
	OperationIndex string `db:"operation_index"`
	Origin         string `db:"origin,omitempty"`
	Tx             string `db:"tx,omitempty"`
	EventId        string `db:"event_id,omitempty"`
	FromChain      string `db:"from_chain,omitempty"`
	ToChain        string `db:"to_chain,omitempty"`
	Receiver       string `db:"receiver,omitempty"`
	Amount         string `db:"amount,omitempty"`
	BundleData     string `db:"bundle_data,omitempty"`
	BundleSalt     string `db:"bundle_salt,omitempty"`
	TokenIndex     string `db:"token_index,omitempty"`
}

func NewTransferRow(
	operationIndex,
	origin,
	tx,
	eventId,
	fromChain,
	toChain,
	receiver,
	amount,
	bundleData,
	bundleSalt,
	tokenIndex string,
) TransferRow {
	return TransferRow{
		OperationIndex: operationIndex,
		Origin:         origin,
		Tx:             tx,
		EventId:        eventId,
		FromChain:      fromChain,
		ToChain:        toChain,
		Receiver:       receiver,
		Amount:         amount,
		BundleData:     bundleData,
		BundleSalt:     bundleSalt,
		TokenIndex:     tokenIndex,
	}
}

type ChangePartiesRow struct {
	OperationIndex string   `db:"operation_index"`
	Parties        []string `db:"parties"`
	NewPublicKey   string   `db:"new_public_key"`
	Signature      string   `db:"signature"`
}

func NewChangePartiesRow(
	operationIndex string,
	parties []string,
	newPublicKey,
	signature string,
) ChangePartiesRow {
	return ChangePartiesRow{
		OperationIndex: operationIndex,
		Parties:        parties,
		NewPublicKey:   newPublicKey,
		Signature:      signature,
	}
}

type ConfirmationRow struct {
	Root           string   `db:"root"`
	Indexes        []string `db:"indexes"`
	SignatureECDSA string   `db:"signature_ecdsa"`
	Creator        string   `db:"creator"`
}

func NewConfirmationRow(
	root string,
	indexes []string,
	signatureECDSA string,
	creator string,
) ConfirmationRow {
	return ConfirmationRow{
		Root:           root,
		Indexes:        indexes,
		SignatureECDSA: signatureECDSA,
		Creator:        creator,
	}
}

package types

// PartyRow represents a single row of the "parties" table
type PartyRow struct {
	Account  string `db:"account"`
	PubKey   string `db:"pub_key"`
	Address  string `db:"address"`
	Verified bool   `db:"verified"`
}

// RarimoCoreParamsRow represents a single row of the "rarimocore_params" table
type RarimoCoreParamsRow struct {
	OneRowID                  bool     `db:"one_row_id"`
	KeyECDSA                  string   `db:"key_ecdsa"`
	Threshold                 uint64   `db:"threshold"`
	IsUpdateRequired          bool     `db:"is_update_required"`
	LastSignature             string   `db:"last_signature"`
	Parties                   []string `db:"parties"`
	AvailableResignBlockDelta uint64   `db:"available_resign_block_delta"`
	Height                    int64    `db:"height"`
}

// OperationRow represents a single row of the "operation" table
type OperationRow struct {
	Index         string `db:"index"`
	OperationType int32  `db:"operation_type"`
	Signed        bool   `db:"signed"`
	Approved      bool   `db:"approved"`
	Creator       string `db:"creator"`
	Timestamp     uint64 `db:"timestamp"`
}

// TransferRow represents a single row of the "transfer" table
type TransferRow struct {
	OperationIndex string            `db:"operation_index"`
	Origin         string            `db:"origin"`
	Tx             string            `db:"tx"`
	EventId        string            `db:"event_id"`
	Receiver       string            `db:"receiver"`
	Amount         string            `db:"amount"`
	BundleData     string            `db:"bundle_data"`
	BundleSalt     string            `db:"bundle_salt"`
	From           *OnChainItemIndex `db:"from"`
	To             *OnChainItemIndex `db:"to"`
	ItemMeta       *ItemMetadata     `db:"item_meta"`
}

// ChangePartiesRow represents a single row of the "change_parties" table
type ChangePartiesRow struct {
	OperationIndex string   `db:"operation_index"`
	Parties        []string `db:"parties"`
	NewPublicKey   string   `db:"new_public_key"`
	Signature      string   `db:"signature"`
}

// ConfirmationRow represents a single row of the "confirmation" table
type ConfirmationRow struct {
	Root           string   `db:"root"`
	Indexes        []string `db:"indexes"`
	SignatureECDSA string   `db:"signature_ecdsa"`
	Creator        string   `db:"creator"`
}

// RarimoCoreVoteRow represents a single row of the "vote" table
type RarimoCoreVoteRow struct {
	Operation string `db:"operation"`
	Validator string `db:"validator"`
	Vote      int32  `db:"vote"`
}

package types

// PartyRow represents a single row of the "parties" table
type PartyRow struct {
	Account         string `db:"account"`
	PubKey          string `db:"pub_key"`
	Address         string `db:"address"`
	Status          int32  `db:"status"`
	ViolationsCount uint64 `db:"violations_count"`
	FreezeEndBlock  uint64 `db:"freeze_end_block"`
	Delegator       string `db:"delegator"`
}

// RarimoCoreParamsRow represents a single row of the "rarimocore_params" table
type RarimoCoreParamsRow struct {
	OneRowID           bool     `db:"one_row_id"`
	KeyECDSA           string   `db:"key_ecdsa"`
	Threshold          uint64   `db:"threshold"`
	IsUpdateRequired   bool     `db:"is_update_required"`
	LastSignature      string   `db:"last_signature"`
	StakeAmount        string   `db:"stake_amount"`
	StakeDenom         string   `db:"stake_denom"`
	MaxViolationsCount uint64   `db:"max_violations_count"`
	FreezeBlocksPeriod uint64   `db:"freeze_blocks_period"`
	Parties            []string `db:"parties"`
	Height             int64    `db:"height"`
}

// OperationRow represents a single row of the "operation" table
type OperationRow struct {
	Index         string `db:"index"`
	OperationType int32  `db:"operation_type"`
	Status        int32  `db:"status"`
	Creator       string `db:"creator"`
	Timestamp     uint64 `db:"timestamp"`
}

// TransferRow represents a single row of the "transfer" table
type TransferRow struct {
	OperationIndex string `db:"operation_index"`
	Origin         string `db:"origin"`
	Tx             string `db:"tx"`
	EventId        string `db:"event_id"`
	Receiver       string `db:"receiver"`
	Amount         string `db:"amount"`
	BundleData     string `db:"bundle_data"`
	BundleSalt     string `db:"bundle_salt"`
	From           string `db:"from"`
	To             string `db:"to"`
	ItemMeta       string `db:"item_meta"`
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

type ViolationReportRow struct {
	Index         string `db:"index"`
	SessionId     string `db:"session_id"`
	Offender      string `db:"offender"`
	ViolationType int32  `db:"violation_type"`
	Sender        string `db:"sender"`
	Msg           string `db:"msg"`
}

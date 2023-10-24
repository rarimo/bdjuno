package types

import (
	rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/rarimo/rarimo-core/x/tokenmanager/types"
)

//--------------------------------------------------------

// Party contains the data of the x/rarimocore module signer instance
type Party struct {
	Account                  string                      `json:"account,omitempty" yaml:"account,omitempty"`
	PubKey                   string                      `json:"pub_key,omitempty" yaml:"pub_key,omitempty"`
	Address                  string                      `json:"address,omitempty" yaml:"address,omitempty"`
	Status                   rarimocoretypes.PartyStatus `json:"status,omitempty" yaml:"status,omitempty"`
	ViolationsCount          uint64                      `json:"violations_count,omitempty" yaml:"violations_count,omitempty"`
	FreezeEndBlock           uint64                      `json:"freeze_end_block,omitempty" yaml:"freeze_end_block,omitempty"`
	Delegator                string                      `json:"delegator,omitempty" yaml:"delegator,omitempty"`
	CommittedGlobalPublicKey string                      `json:"committed_global_public_key,omitempty"`
	ReportedSessions         []string                    `json:"reported_sessions,omitempty"`
}

// NewParty allows to build a new Party
func NewParty(p rarimocoretypes.Party) Party {
	return Party{
		Account:                  p.Account,
		PubKey:                   p.PubKey,
		Address:                  p.Address,
		Status:                   p.Status,
		ViolationsCount:          p.ViolationsCount,
		FreezeEndBlock:           p.FreezeEndBlock,
		Delegator:                p.Delegator,
		CommittedGlobalPublicKey: p.CommittedGlobalPublicKey,
		ReportedSessions:         p.ReportedSessions,
	}
}

//--------------------------------------------------------

// RarimoCoreParams contains the data of the x/rarimocore module params instance
type RarimoCoreParams struct {
	KeyECDSA           string   `json:"key_ecdsa,omitempty" yaml:"key_ecdsa,omitempty"`
	Threshold          uint64   `json:"threshold,omitempty" yaml:"threshold,omitempty"`
	IsUpdateRequired   bool     `json:"is_update_required,omitempty" yaml:"is_update_required,omitempty"`
	LastSignature      string   `json:"last_signature,omitempty" yaml:"last_signature,omitempty"`
	StakeAmount        string   `json:"stake_amount,omitempty" yaml:"stake_amount,omitempty"`
	StakeDenom         string   `json:"stake_denom,omitempty" yaml:"stake_denom,omitempty"`
	MaxViolationsCount uint64   `json:"max_violations_count,omitempty" yaml:"max_violations_count,omitempty"`
	FreezeBlocksPeriod uint64   `json:"freeze_blocks_period,omitempty" yaml:"freeze_blocks_period,omitempty"`
	Parties            []string `json:"parties,omitempty" yaml:"parties,omitempty"`
	Height             int64    `json:"height,omitempty" yaml:"height,omitempty"`
}

// NewRarimoCoreParams allows to build a new RarimoCoreParams instance
func NewRarimoCoreParams(p rarimocoretypes.Params, height int64) *RarimoCoreParams {
	parties := make([]string, len(p.Parties))
	for i, party := range p.Parties {
		parties[i] = party.Account
	}
	return &RarimoCoreParams{
		KeyECDSA:           p.KeyECDSA,
		Threshold:          p.Threshold,
		IsUpdateRequired:   p.IsUpdateRequired,
		LastSignature:      p.LastSignature,
		StakeAmount:        p.StakeAmount,
		StakeDenom:         p.StakeDenom,
		MaxViolationsCount: p.MaxViolationsCount,
		FreezeBlocksPeriod: p.FreezeBlocksPeriod,
		Parties:            parties,
		Height:             height,
	}
}

//--------------------------------------------------------

// Operation represents a single operation instance
type Operation struct {
	Index         string                   `json:"index,omitempty" yaml:"index,omitempty"`
	OperationType rarimocoretypes.OpType   `json:"operation_type" yaml:"operation_type"`
	Status        rarimocoretypes.OpStatus `json:"status" yaml:"status"`
	Creator       string                   `json:"creator,omitempty" yaml:"creator,omitempty"`
	Timestamp     int64                    `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
}

// NewOperation allows to build a new Operation instance
func NewOperation(index string, opType, status int32, creator string, timestamp uint64) Operation {
	return Operation{
		Index:         index,
		OperationType: rarimocoretypes.OpType(opType),
		Status:        rarimocoretypes.OpStatus(status),
		Creator:       creator,
		Timestamp:     int64(timestamp),
	}
}

// OperationFromCore allows to build a new Operation instance from a rarimocoretypes.Operation instance
func OperationFromCore(operation rarimocoretypes.Operation) Operation {
	return Operation{
		Index:         operation.Index,
		OperationType: operation.OperationType,
		Status:        operation.Status,
		Creator:       operation.Creator,
		Timestamp:     int64(operation.Timestamp),
	}
}

//--------------------------------------------------------

// Transfer represents a single transfer instance
type Transfer struct {
	OperationIndex string            `json:"operation_index,omitempty" yaml:"operation_index,omitempty"`
	Origin         string            `json:"origin,omitempty" yaml:"origin,omitempty"`
	Tx             string            `json:"tx,omitempty" yaml:"tx,omitempty"`
	EventID        string            `json:"event_id,omitempty" yaml:"event_id,omitempty"`
	Receiver       string            `json:"receiver,omitempty" yaml:"receiver,omitempty"`
	Amount         string            `json:"amount,omitempty" yaml:"amount,omitempty"`
	BundleData     string            `json:"bundle_data,omitempty" yaml:"bundle_data,omitempty"`
	BundleSalt     string            `json:"bundle_salt,omitempty" yaml:"bundle_salt,omitempty"`
	From           *OnChainItemIndex `json:"from,omitempty" yaml:"from,omitempty"`
	To             *OnChainItemIndex `json:"to,omitempty" yaml:"to,omitempty"`
	ItemMeta       *ItemMetadata     `json:"item_meta,omitempty" yaml:"item_meta,omitempty"`
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
		ItemMeta:       ItemMetadataFromCore(t.Meta),
		From:           OnChainItemIndexFromCore(&t.From),
		To:             OnChainItemIndexFromCore(&t.To),
	}
}

//--------------------------------------------------------

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

//--------------------------------------------------------

type ContractUpgrade struct {
	OperationIndex            string            `json:"operation_index,omitempty" yaml:"operation_index,omitempty"`
	TargetContract            string            `json:"target_contract,omitempty" yaml:"target_contract,omitempty"`
	Chain                     string            `json:"chain,omitempty" yaml:"chain,omitempty"`
	NewImplementationContract string            `json:"new_implementation_contract,omitempty" yaml:"new_implementation_contract,omitempty"`
	Hash                      string            `json:"hash,omitempty" yaml:"hash,omitempty"`
	BufferAccount             string            `json:"buffer_account,omitempty" yaml:"buffer_account,omitempty"`
	Nonce                     string            `json:"nonce,omitempty" yaml:"nonce,omitempty"`
	Type                      types.UpgradeType `json:"type,omitempty" yaml:"type,omitempty"`
}

// NewContractUpdate allows to build a new ContractUpgrade instance
func NewContractUpdate(operationIndex string, c rarimocoretypes.ContractUpgrade) ContractUpgrade {
	return ContractUpgrade{
		OperationIndex:            operationIndex,
		TargetContract:            c.TargetContract,
		Chain:                     c.Chain,
		NewImplementationContract: c.NewImplementationContract,
		Hash:                      c.Hash,
		BufferAccount:             c.BufferAccount,
		Nonce:                     c.Nonce,
		Type:                      c.Type,
	}
}

//--------------------------------------------------------

type FeeTokenManagement struct {
	OperationIndex   string                                 `json:"operation_index,omitempty" yaml:"operation_index,omitempty"`
	OpType           rarimocoretypes.FeeTokenManagementType `json:"op_type,omitempty" yaml:"op_type,omitempty"`
	FeeTokenContract string                                 `json:"fee_token_contract" yaml:"fee_token_contract"`
	FeeTokenAmount   string                                 `json:"fee_token_amount" yaml:"fee_token_amount"`
	Chain            string                                 `json:"chain,omitempty" yaml:"chain,omitempty"`
	Receiver         string                                 `json:"receiver,omitempty" yaml:"receiver,omitempty"`
	Nonce            string                                 `json:"nonce,omitempty" yaml:"nonce,omitempty"`
}

// NewFeeTokenManagement allows to build a new FeeTokenManagement instance
func NewFeeTokenManagement(operationIndex string, f rarimocoretypes.FeeTokenManagement) FeeTokenManagement {
	return FeeTokenManagement{
		OperationIndex:   operationIndex,
		OpType:           f.OpType,
		FeeTokenContract: f.Token.Contract,
		FeeTokenAmount:   f.Token.Amount,
		Chain:            f.Chain,
		Receiver:         f.Receiver,
		Nonce:            f.Nonce,
	}
}

//--------------------------------------------------------

type IdentityDefaultTransfer struct {
	OperationIndex          string `json:"operation_index,omitempty" yaml:"operation_index,omitempty"`
	Contract                string `json:"contract,omitempty" yaml:"contract,omitempty"`
	Chain                   string `json:"chain,omitempty" yaml:"chain,omitempty"`
	GISTHash                string `json:"gisthash,omitempty" yaml:"gisthash,omitempty"`
	Id                      string `json:"id,omitempty" yaml:"id,omitempty"`
	StateHash               string `json:"state_hash,omitempty" yaml:"state_hash,omitempty"`
	StateCreatedAtTimestamp string `json:"state_created_at_timestamp,omitempty" yaml:"state_created_at_timestamp,omitempty"`
	StateCreatedAtBlock     string `json:"state_created_at_block,omitempty" yaml:"state_created_at_block,omitempty"`
	StateReplacedBy         string `json:"state_replaced_by,omitempty" yaml:"state_replaced_by,omitempty"`
	GISTReplacedBy          string `json:"gistreplaced_by,omitempty" yaml:"gistreplaced_by,omitempty"`
	GISTCreatedAtTimestamp  string `json:"gistcreated_at_timestamp,omitempty" yaml:"gistcreated_at_timestamp,omitempty"`
	GISTCreatedAtBlock      string `json:"gistcreated_at_block,omitempty" yaml:"gistcreated_at_block,omitempty"`
	ReplacedStateHash       string `json:"replaced_state_hash,omitempty" yaml:"replaced_state_hash,omitempty"`
	ReplacedGISTHash        string `json:"replaced_gist_hash,omitempty" yaml:"replaced_gist_hash,omitempty"`
}

// NewIdentityDefaultTransfer allows to build a new IdentityDefaultTransfer instance
func NewIdentityDefaultTransfer(operationIndex string, i rarimocoretypes.IdentityDefaultTransfer) IdentityDefaultTransfer {
	return IdentityDefaultTransfer{
		OperationIndex:          operationIndex,
		Contract:                i.Contract,
		Chain:                   i.Chain,
		GISTHash:                i.GISTHash,
		Id:                      i.Id,
		StateHash:               i.StateHash,
		StateCreatedAtTimestamp: i.StateCreatedAtTimestamp,
		StateCreatedAtBlock:     i.StateCreatedAtBlock,
		StateReplacedBy:         i.StateReplacedBy,
		GISTReplacedBy:          i.GISTReplacedBy,
		GISTCreatedAtTimestamp:  i.GISTCreatedAtTimestamp,
		GISTCreatedAtBlock:      i.GISTCreatedAtBlock,
		ReplacedStateHash:       i.ReplacedStateHash,
		ReplacedGISTHash:        i.ReplacedGISTHash,
	}

}

//--------------------------------------------------------

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

//--------------------------------------------------------

// RarimoCoreVote represents a single vote instance
type RarimoCoreVote struct {
	Operation string                   `json:"operation,omitempty" yaml:"operation,omitempty"`
	Validator string                   `json:"validator,omitempty" yaml:"validator,omitempty"`
	Vote      rarimocoretypes.VoteType `json:"vote" yaml:"vote"`
}

// NewRarimoCoreVote allows to build a new RarimoCoreVote instance
func NewRarimoCoreVote(operation, validator string, vote int32) RarimoCoreVote {
	return RarimoCoreVote{
		Operation: operation,
		Validator: validator,
		Vote:      rarimocoretypes.VoteType(vote),
	}
}

// RarimoCoreVoteFromCore allows to build a new RarimoCoreVote instance from a rarimocoretypes.Vote instance
func RarimoCoreVoteFromCore(vote rarimocoretypes.Vote) RarimoCoreVote {
	return RarimoCoreVote{
		Operation: vote.Index.Operation,
		Validator: vote.Index.Validator,
		Vote:      vote.Vote,
	}
}

//--------------------------------------------------------

// ViolationReport represents a single violation report instance
type ViolationReport struct {
	Index         string                        `json:"index,omitempty" yaml:"index,omitempty"`
	SessionId     string                        `json:"session_id,omitempty" yaml:"session_id,omitempty"`
	Offender      string                        `json:"offender,omitempty" yaml:"offender,omitempty"`
	ViolationType rarimocoretypes.ViolationType `json:"violation_type,omitempty" yaml:"violation_type,omitempty"`
	Sender        string                        `json:"sender,omitempty" yaml:"sender,omitempty"`
	Msg           string                        `json:"msg,omitempty" yaml:"msg,omitempty"`
}

// ViolationReportFromCore allows to build a new ViolationReport instance from a rarimocoretypes.ViolationReport instance
func ViolationReportFromCore(report rarimocoretypes.ViolationReport) ViolationReport {
	return ViolationReport{
		Index:         string(rarimocoretypes.ViolationReportKey(report.Index)),
		SessionId:     report.Index.SessionId,
		Offender:      report.Index.Offender,
		ViolationType: report.Index.ViolationType,
		Sender:        report.Index.Sender,
		Msg:           report.Msg,
	}
}

//--------------------------------------------------------

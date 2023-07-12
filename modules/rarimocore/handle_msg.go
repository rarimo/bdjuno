package rarimocore

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	juno "github.com/forbole/juno/v4/types"
	"gitlab.com/rarimo/bdjuno/types"
	oracletypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"math/big"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *oracletypes.MsgCreateTransferOp:
		return m.handleMsgCreateTransferOp(tx, cosmosMsg)
	case *rarimocoretypes.MsgCreateViolationReport:
		return m.handleMsgCreateViolationReport(tx, cosmosMsg)
	case *rarimocoretypes.MsgCreateConfirmation:
		return m.handleMsgCreateConfirmation(tx, cosmosMsg)
	case *rarimocoretypes.MsgCreateChangePartiesOp:
		return m.handleMsgCreateChangePartiesOp(tx, cosmosMsg)
	case *oracletypes.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)
	case *rarimocoretypes.MsgStake, *rarimocoretypes.MsgUnstake, *rarimocoretypes.MsgSetupInitial, *rarimocoretypes.MsgChangePartyAddress:
		return m.UpdateParams(tx.Height)
	}

	return nil
}

func (m *Module) handleMsgCreateViolationReport(tx *juno.Tx, msg *rarimocoretypes.MsgCreateViolationReport) error {
	rawReport, err := m.source.ViolationReport(tx.Height, msg.SessionId, msg.Offender, msg.SessionId, msg.ViolationType)
	if err != nil {
		return fmt.Errorf("failed to get violation report: %s", err)
	}

	report := types.ViolationReportFromCore(rawReport)

	err = m.db.SaveViolationReports([]types.ViolationReport{report})
	if err != nil {
		return fmt.Errorf("failed to save violation report: %s", err)
	}

	err = m.UpdateParams(tx.Height)
	if err != nil {
		return fmt.Errorf("failed to update last rarimocore params: %s", err)
	}

	return nil
}

func (m *Module) handleMsgVote(tx *juno.Tx, msg *oracletypes.MsgVote) error {
	rawOp, err := m.source.Operation(tx.Height, msg.Operation)
	if err != nil {
		return fmt.Errorf("failed to get change operation: %s", err)
	}

	op := types.OperationFromCore(rawOp)

	err = m.db.SaveRarimoCoreVotes(
		[]types.RarimoCoreVote{types.NewRarimoCoreVote(msg.Operation, msg.Index.Account, int32(msg.Vote))},
	)
	if err != nil {
		return fmt.Errorf("failed to save vote: %s", err)
	}

	if op.OperationType == rarimocoretypes.OpType_TRANSFER && op.Status == rarimocoretypes.OpStatus_APPROVED {
		err = m.handleApproveTransfer(tx, rawOp)
	}

	err = m.db.UpdateOperation(op)
	if err != nil {
		return fmt.Errorf("failed to update operation: %s", err)
	}
	return nil

}

func (m *Module) handleApproveTransfer(tx *juno.Tx, op rarimocoretypes.Operation) error {
	transfer, err := getTransfer(op)
	if err != nil {
		return fmt.Errorf("failed to get transfer: %s", err)
	}

	if err != nil {
		return fmt.Errorf("failed to get collection data: %s", err)
	}

	from, err := m.tokenmanagerSource.OnChainItem(tx.Height, *transfer.From)
	if err != nil {
		return fmt.Errorf("failed to get on chain item: %s", err)
	}

	item, err := m.tokenmanagerSource.Item(tx.Height, transfer.Origin)
	if err != nil {
		return fmt.Errorf("failed to get item: %s", err)
	}

	err = m.db.UpsertItem(types.ItemFromCore(item))
	if err != nil {
		return fmt.Errorf("failed to save item: %s", err)
	}

	if item.Meta.Seed != "" {
		seed, err := m.tokenmanagerSource.Seed(tx.Height, item.Meta.Seed)
		if err != nil {
			return fmt.Errorf("failed to get seed: %s", err)
		}

		err = m.db.UpsertSeed(types.SeedFromCore(seed))
		if err != nil {
			return fmt.Errorf("failed to save seed: %s", err)
		}
	}

	to, err := m.tokenmanagerSource.OnChainItem(tx.Height, *transfer.To)
	if err != nil {
		return fmt.Errorf("failed to get on chain item: %s", err)
	}

	err = m.db.SaveOnChainItems([]types.OnChainItem{
		types.OnChainItemFromCore(from),
		types.OnChainItemFromCore(to),
	})
	if err != nil {
		return fmt.Errorf("failed to save on chain item: %s", err)
	}

	return nil

}

func (m *Module) handleMsgCreateTransferOp(tx *juno.Tx, msg *oracletypes.MsgCreateTransferOp) error {
	index := hexutil.Encode(crypto.Keccak256(
		[]byte(msg.Tx),
		[]byte(msg.EventId),
		[]byte(msg.From.Chain),
	))

	op, err := m.source.Operation(tx.Height, index)
	if err != nil {
		return fmt.Errorf("failed to get transfer operation: %s", err)
	}

	err = m.saveOperations([]rarimocoretypes.Operation{op})
	if err != nil {
		return fmt.Errorf("failed to save transfer operation: %s", err)
	}

	if op.Status == rarimocoretypes.OpStatus_INITIALIZED || op.Status == rarimocoretypes.OpStatus_APPROVED {
		return nil
	}

	err = m.db.RemoveRarimoCoreVotes(op.Index)
	if err != nil {
		return fmt.Errorf("failed to remove votes: %s", err)
	}

	return nil
}

func (m *Module) handleMsgCreateChangePartiesOp(tx *juno.Tx, msg *rarimocoretypes.MsgCreateChangePartiesOp) error {
	var changeOp = &rarimocoretypes.ChangeParties{
		Parties:      msg.NewSet,
		Signature:    msg.Signature,
		NewPublicKey: msg.NewPublicKey,
	}

	content, _ := pkg.GetChangePartiesContent(changeOp)

	op, err := m.source.Operation(tx.Height, hexutil.Encode(content.CalculateHash()))
	if err != nil {
		return fmt.Errorf("failed to get change parties operation: %s", err)
	}

	err = m.saveOperations([]rarimocoretypes.Operation{op})
	if err != nil {
		return fmt.Errorf("failed to save change parties operation: %s", err)
	}

	return nil
}

func (m *Module) handleMsgCreateConfirmation(tx *juno.Tx, msg *rarimocoretypes.MsgCreateConfirmation) error {
	var confirmation = rarimocoretypes.Confirmation{
		Creator:        msg.Creator,
		Root:           msg.Root,
		Indexes:        msg.Indexes,
		SignatureECDSA: msg.SignatureECDSA,
	}

	ops := make([]rarimocoretypes.Operation, len(msg.Indexes))

	for i, index := range msg.Indexes {
		op, err := m.source.Operation(tx.Height, index)
		if err != nil {
			return fmt.Errorf("failed to get operation: %s", err)
		}

		ops[i] = op
	}

	err := m.updateOperations(ops)
	if err != nil {
		return fmt.Errorf("failed to update operations: %s", err)
	}

	err = m.saveConfirmations([]rarimocoretypes.Confirmation{confirmation})
	if err != nil {
		return fmt.Errorf("failed to save confirmation: %s", err)
	}

	err = m.UpdateParams(tx.Height)
	if err != nil {
		return fmt.Errorf("failed to update last rarimocore params: %s", err)
	}

	return nil
}

func (m *Module) HandleUpdateContract(height int64, details tokenmanagertypes.ContractUpgradeDetails) error {
	params, err := m.tokenmanagerSource.Params(height)
	if err != nil {
		return fmt.Errorf("failed to tokenmanager get params: %s", err)
	}

	network, ok := GetNetwork(params, details.Chain)
	if !ok {
		return fmt.Errorf("failed to get network")
	}

	upgrade := &rarimocoretypes.ContractUpgrade{
		TargetContract:            details.TargetContract,
		Chain:                     details.Chain,
		NewImplementationContract: details.NewImplementationContract,
		Hash:                      details.Hash,
		BufferAccount:             details.BufferAccount,
		Nonce:                     details.Nonce,
		Type:                      details.Type,
	}

	content, err := pkg.GetContractUpgradeContent(network, upgrade)
	if err != nil {
		return fmt.Errorf("error creating content %s", err)
	}

	index := hexutil.Encode(crypto.Keccak256(big.NewInt(height).Bytes(), content.CalculateHash()))

	op, err := m.source.Operation(height, index)
	if err != nil {
		return fmt.Errorf("failed to get contract upgrade operation: %s", err)
	}

	err = m.saveOperations([]rarimocoretypes.Operation{op})
	if err != nil {
		return fmt.Errorf("failed to save contract upgrade operation: %s", err)
	}

	return nil
}

func (m *Module) GetFeeToken(height int64, chain, contract string) (*tokenmanagertypes.FeeToken, error) {
	params, err := m.tokenmanagerSource.Params(height)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenmanager get params: %s", err)
	}

	network, ok := GetNetwork(params, chain)
	if !ok {
		return nil, fmt.Errorf("failed to get network")
	}

	feeparams := network.GetFeeParams()
	if feeparams == nil {
		return nil, fmt.Errorf("failed to get fee params")
	}

	feetoken := feeparams.GetFeeToken(contract)
	if feetoken == nil {
		return nil, fmt.Errorf("failed to get fee token")
	}

	return feetoken, nil
}

func GetNetwork(params tokenmanagertypes.Params, name string) (param tokenmanagertypes.Network, ok bool) {
	for _, network := range params.Networks {
		if network.Name == name {
			return *network, true
		}
	}

	return
}

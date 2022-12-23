package rarimocore

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	juno "github.com/forbole/juno/v3/types"
	"gitlab.com/rarimo/bdjuno/types"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/operation/origin"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *rarimocoretypes.MsgCreateTransferOp:
		return m.handleMsgCreateTransferOp(tx, cosmosMsg)

	case *rarimocoretypes.MsgCreateChangePartiesOp:
		return m.handleMsgCreateChangePartiesOp(tx, cosmosMsg)

	case *rarimocoretypes.MsgCreateConfirmation:
		return m.handleCreateConfirmation(tx, cosmosMsg)

		//case *rarimocoretypes.MsgSetupInitial:
		//	return m.handleMsgVote(tx, cosmosMsg)
	}

	return nil
}

func (m *Module) handleMsgCreateTransferOp(tx *juno.Tx, msg *rarimocoretypes.MsgCreateTransferOp) error {
	or := origin.NewDefaultOriginBuilder().
		SetTxHash(msg.Tx).
		SetOpId(msg.EventId).
		SetCurrentNetwork(msg.FromChain).
		Build().
		GetOrigin()

	index := hexutil.Encode(or[:])

	op, err := m.source.Operation(tx.Height, index)
	if err != nil {
		return fmt.Errorf("failed to get operation from source: %s", err)
	}

	err = m.saveOperations([]rarimocoretypes.Operation{op})
	if err != nil {
		return fmt.Errorf("failed to save operation: %s", err)
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
	index := hexutil.Encode(content.CalculateHash())

	op, err := m.source.Operation(tx.Height, index)
	if err != nil {
		return fmt.Errorf("failed to get operation: %s", err)
	}

	err = m.saveOperations([]rarimocoretypes.Operation{op})
	if err != nil {
		return fmt.Errorf("failed to save operation: %s", err)
	}

	return nil
}

func (m *Module) handleCreateConfirmation(tx *juno.Tx, msg *rarimocoretypes.MsgCreateConfirmation) error {
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

	params, err := m.source.Params(tx.Height)
	if err != nil {
		return fmt.Errorf("failed to get params: %s", err)
	}

	err = m.db.SaveRarimoCoreParams(types.NewRarimoCoreParams(params, tx.Height))
	if err != nil {
		return fmt.Errorf("failed to update last rarimocore params: %s", err)
	}

	return nil
}

func (m *Module) updateOperations(slice []rarimocoretypes.Operation) error {
	operations := coreOperationsToInternal(slice)

	// Update the operations
	for _, operation := range operations {
		err := m.db.UpdateOperation(operation)
		if err != nil {
			return err
		}
	}

	_, changeParties, err := m.getOperationDetails(slice)
	if err != nil {
		return nil
	}

	// Update the change parties
	for _, changeParty := range changeParties {
		err = m.db.UpdateChangeParties(changeParty)
		if err != nil {
			return err
		}
	}

	return nil
}

func coreOperationsToInternal(slice []rarimocoretypes.Operation) []types.Operation {
	operations := make([]types.Operation, len(slice))
	for i, operation := range slice {
		operations[i] = types.NewOperation(
			operation.Index,
			operation.OperationType,
			operation.Signed,
			operation.Creator,
			operation.Timestamp,
		)
	}

	return operations
}

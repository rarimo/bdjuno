package rarimocore

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	juno "github.com/forbole/juno/v3/types"
	"gitlab.com/rarimo/bdjuno/types"
	"gitlab.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	"math/big"
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
		return m.handleMsgCreateConfirmation(tx, cosmosMsg)
	case *rarimocoretypes.MsgResignOperation:
		return m.handleMsgResignOperation(tx, cosmosMsg)
	case *rarimocoretypes.MsgSetupInitial, *rarimocoretypes.MsgChangePartyAddress:
		return m.UpdateParams(tx.Height)
	}

	return nil
}

func (m *Module) handleMsgResignOperation(tx *juno.Tx, msg *rarimocoretypes.MsgResignOperation) error {
	op, err := m.source.Operation(tx.Height, msg.ChangeOpIndex)
	if err != nil {
		return fmt.Errorf("failed to get change operation: %s", err)
	}

	err = m.db.UpdateOperation(types.OperationFromCore(op))
	if err != nil {
		return fmt.Errorf("failed to update operation: %s", err)
	}

	return nil
}

func (m *Module) handleMsgCreateTransferOp(tx *juno.Tx, msg *rarimocoretypes.MsgCreateTransferOp) error {
	index := hexutil.Encode(crypto.Keccak256([]byte(msg.Tx),
		[]byte(msg.EventId),
		[]byte(msg.From.Chain),
		big.NewInt(tx.Height).Bytes(),
	))

	err := m.handleNewOperation(tx.Height, index)
	if err != nil {
		return fmt.Errorf("failed to handle new create transfer operation: %s", err)
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

	err := m.handleNewOperation(tx.Height, hexutil.Encode(content.CalculateHash()))
	if err != nil {
		return fmt.Errorf("failed to handle new create change parties operation: %s", err)
	}

	return nil
}

func (m *Module) handleNewOperation(height int64, index string) error {
	op, err := m.source.Operation(height, index)
	if err != nil {
		return fmt.Errorf("failed to get operation: %s", err)
	}

	return m.db.Transaction(func() error {
		err = m.saveOperations([]rarimocoretypes.Operation{op})
		if err != nil {
			return fmt.Errorf("failed to save operation: %s", err)
		}

		err = m.UpdateParams(height)
		if err != nil {
			return fmt.Errorf("failed to update last rarimocore params: %s", err)
		}

		return nil
	})
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

	return m.db.Transaction(func() error {
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
	})
}

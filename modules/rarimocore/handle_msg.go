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
	case *rarimocoretypes.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)
	case *rarimocoretypes.MsgSetupInitial, *rarimocoretypes.MsgChangePartyAddress:
		return m.UpdateParams(tx.Height)
	}

	return nil
}

func (m *Module) handleMsgVote(tx *juno.Tx, msg *rarimocoretypes.MsgVote) error {
	rawOp, err := m.source.Operation(tx.Height, msg.Operation)
	if err != nil {
		return fmt.Errorf("failed to get change operation: %s", err)
	}

	op := types.OperationFromCore(rawOp)

	return m.db.Transaction(func() error {
		err = m.db.SaveRarimoCoreVotes(
			[]types.RarimoCoreVote{types.NewRarimoCoreVote(msg.Operation, msg.Creator, int32(msg.Vote))},
		)

		if op.Status != rarimocoretypes.OpStatus_INITIALIZED {
			return nil
		}

		if op.OperationType == rarimocoretypes.OpType_TRANSFER && op.Status == rarimocoretypes.OpStatus_APPROVED {
			err = m.handleApproveTransfer(tx, rawOp)
		}

		err = m.db.UpdateOperation(op)
		if err != nil {
			return fmt.Errorf("failed to update operation: %s", err)
		}
		return nil
	})
}

func (m *Module) handleApproveTransfer(tx *juno.Tx, op rarimocoretypes.Operation) error {
	transfer, err := getTransfer(op)
	if err != nil {
		return fmt.Errorf("failed to get transfer: %s", err)
	}

	if err != nil {
		return fmt.Errorf("failed to get collection data: %s", err)
	}

	from, err := m.tokenManagerSource.OnChainItem(tx.Height, *transfer.From)
	if err != nil {
		return fmt.Errorf("failed to get on chain item: %s", err)
	}

	item, err := m.tokenManagerSource.Item(tx.Height, transfer.Origin)
	if err != nil {
		return fmt.Errorf("failed to get item: %s", err)
	}

	return m.db.Transaction(func() error {
		err = m.db.UpsertItem(types.ItemFromCore(item))
		if err != nil {
			return fmt.Errorf("failed to save item: %s", err)
		}

		if item.Meta.Seed != "" {
			seed, err := m.tokenManagerSource.Seed(tx.Height, item.Meta.Seed)
			if err != nil {
				return fmt.Errorf("failed to get seed: %s", err)
			}

			err = m.db.UpsertSeed(types.SeedFromCore(seed))
			if err != nil {
				return fmt.Errorf("failed to save seed: %s", err)
			}
		}

		to, err := m.tokenManagerSource.OnChainItem(tx.Height, *transfer.To)
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
	})
}

func (m *Module) handleMsgCreateTransferOp(tx *juno.Tx, msg *rarimocoretypes.MsgCreateTransferOp) error {
	index := hexutil.Encode(crypto.Keccak256([]byte(msg.Tx),
		[]byte(msg.EventId),
		[]byte(msg.From.Chain),
		big.NewInt(tx.Height).Bytes(),
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

package oraclemanager

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	juno "github.com/forbole/juno/v4/types"
	"github.com/rarimo/bdjuno/types"
	oracletypes "github.com/rarimo/rarimo-core/x/oraclemanager/types"
	"github.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *oracletypes.MsgStake:
		return m.HandleOracle(tx.Height, cosmosMsg.Index.Chain, cosmosMsg.Index.Account)
	case *oracletypes.MsgUnstake:
		return m.HandleOracle(tx.Height, cosmosMsg.Index.Chain, cosmosMsg.Index.Account)
	case *oracletypes.MsgUnjail:
		return m.HandleOracle(tx.Height, cosmosMsg.Index.Chain, cosmosMsg.Index.Account)
	case *oracletypes.MsgCreateTransferOp:
		return m.HandleOracle(tx.Height, cosmosMsg.From.Chain, cosmosMsg.Creator)
	case *oracletypes.MsgVote:
		return m.HandleOracle(tx.Height, cosmosMsg.Index.Chain, cosmosMsg.Index.Account)
	case *oracletypes.MsgCreateIdentityDefaultTransferOp:
		return m.handleCreateIdentityDefault(tx.Height, cosmosMsg)
	case *oracletypes.MsgCreateIdentityGISTTransferOp:
		return m.handleCreateIdentityGIST(tx.Height, cosmosMsg)
	case *oracletypes.MsgCreateIdentityStateTransferOp:
		return m.handleCreateIdentityState(tx.Height, cosmosMsg)
	}

	return nil
}

func (m *Module) handleCreateIdentityDefault(height int64, msg *oracletypes.MsgCreateIdentityDefaultTransferOp) error {
	transfer := &rarimocoretypes.IdentityDefaultTransfer{
		Contract:                msg.Contract,
		Chain:                   msg.Chain,
		GISTHash:                msg.GISTHash,
		Id:                      msg.Id,
		StateHash:               msg.StateHash,
		StateCreatedAtTimestamp: msg.StateCreatedAtTimestamp,
		StateCreatedAtBlock:     msg.StateCreatedAtBlock,
		StateReplacedBy:         msg.StateReplacedBy,
		GISTReplacedBy:          msg.GISTReplacedBy,
		GISTCreatedAtTimestamp:  msg.GISTCreatedAtTimestamp,
		GISTCreatedAtBlock:      msg.GISTCreatedAtBlock,
		ReplacedStateHash:       msg.ReplacedStateHash,
		ReplacedGISTHash:        msg.ReplacedGISTtHash,
	}

	content, err := pkg.GetIdentityDefaultTransferContent(transfer)
	if err != nil {
		return fmt.Errorf("error creating content %s", err)
	}

	index := hexutil.Encode(content.CalculateHash())

	err = m.rc.SaveOperationByIndex(height, index)
	if err != nil {
		return fmt.Errorf("error saving operation %s", err)
	}

	return nil
}

func (m *Module) handleCreateIdentityGIST(height int64, msg *oracletypes.MsgCreateIdentityGISTTransferOp) error {
	transfer := &rarimocoretypes.IdentityGISTTransfer{
		Contract:               msg.Contract,
		Chain:                  msg.Chain,
		GISTHash:               msg.GISTHash,
		GISTCreatedAtTimestamp: msg.GISTCreatedAtTimestamp,
		GISTCreatedAtBlock:     msg.GISTCreatedAtBlock,
		ReplacedGISTHash:       msg.ReplacedGISTtHash,
	}

	content, err := pkg.GetIdentityGISTTransferContent(transfer)
	if err != nil {
		return fmt.Errorf("error creating content %s", err)
	}

	index := hexutil.Encode(content.CalculateHash())

	err = m.rc.SaveOperationByIndex(height, index)
	if err != nil {
		return fmt.Errorf("error saving operation %s", err)
	}

	return nil
}

func (m *Module) handleCreateIdentityState(height int64, msg *oracletypes.MsgCreateIdentityStateTransferOp) error {
	transfer := &rarimocoretypes.IdentityStateTransfer{
		Contract:                msg.Contract,
		Chain:                   msg.Chain,
		Id:                      msg.Id,
		StateHash:               msg.StateHash,
		StateCreatedAtTimestamp: msg.StateCreatedAtTimestamp,
		StateCreatedAtBlock:     msg.StateCreatedAtBlock,
		ReplacedStateHash:       msg.ReplacedStateHash,
	}

	content, err := pkg.GetIdentityStateTransferContent(transfer)
	if err != nil {
		return fmt.Errorf("error creating content %s", err)
	}

	index := hexutil.Encode(content.CalculateHash())

	err = m.rc.SaveOperationByIndex(height, index)
	if err != nil {
		return fmt.Errorf("error saving operation %s", err)
	}

	return nil
}

func (m *Module) HandleOracle(height int64, chain, account string) error {
	raw, err := m.source.Oracle(height, chain, account)
	if err != nil {
		return fmt.Errorf("failed to get oracle: %s", err)
	}

	oracle := types.OracleFromCore(raw)

	err = m.db.SaveOracles([]types.Oracle{oracle})
	if err != nil {
		return fmt.Errorf("failed to save oracle: %s", err)
	}

	return nil

}

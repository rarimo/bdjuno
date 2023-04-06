package oraclemanager

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
	"gitlab.com/rarimo/bdjuno/types"
	oracletypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *oracletypes.MsgStake:
		return m.handleMsgStake(tx, cosmosMsg)
	case *oracletypes.MsgUnstake:
		return m.handleMsgUnstake(tx, cosmosMsg)
	case *oracletypes.MsgUnjail:
		return m.handleMsgUnjail(tx, cosmosMsg)
	case *oracletypes.MsgCreateTransferOp:
		return m.handleMsgCreateTransferOp(tx, cosmosMsg)
	case *oracletypes.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)
	}

	return nil
}

func (m *Module) handleMsgStake(tx *juno.Tx, msg *oracletypes.MsgStake) error {
	return m.HandleOracle(tx.Height, msg.Index.Chain, msg.Index.Account)
}

func (m *Module) handleMsgUnstake(tx *juno.Tx, msg *oracletypes.MsgUnstake) error {
	return m.HandleOracle(tx.Height, msg.Index.Chain, msg.Index.Account)
}

func (m *Module) handleMsgUnjail(tx *juno.Tx, msg *oracletypes.MsgUnjail) error {
	return m.HandleOracle(tx.Height, msg.Index.Chain, msg.Index.Account)
}

func (m *Module) handleMsgCreateTransferOp(tx *juno.Tx, msg *oracletypes.MsgCreateTransferOp) error {
	return m.HandleOracle(tx.Height, msg.From.Chain, msg.Creator)
}

func (m *Module) handleMsgVote(tx *juno.Tx, msg *oracletypes.MsgVote) error {
	return m.HandleOracle(tx.Height, msg.Index.Chain, msg.Index.Account)
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

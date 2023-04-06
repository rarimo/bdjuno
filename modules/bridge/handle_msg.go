package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
	"gitlab.com/rarimo/bdjuno/types"
	bridgetypes "gitlab.com/rarimo/rarimo-core/x/bridge/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *bridgetypes.MsgWithdrawNative:
		return m.handleMsgWithdrawNative(tx, cosmosMsg)
	}

	return nil
}

func (m *Module) handleMsgWithdrawNative(tx *juno.Tx, msg *bridgetypes.MsgWithdrawNative) error {
	return m.db.SaveHashes([]types.Hash{
		{
			Index: msg.Origin,
		},
	})
}

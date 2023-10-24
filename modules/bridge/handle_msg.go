package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v4/types"
	"github.com/rarimo/bdjuno/types"
	bridgetypes "github.com/rarimo/rarimo-core/x/bridge/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *bridgetypes.MsgWithdrawNative:
		return m.handleMsgWithdrawNative(cosmosMsg.Origin)
	case *bridgetypes.MsgWithdrawFee:
		return m.handleMsgWithdrawNative(cosmosMsg.Origin)
	}

	return nil
}

func (m *Module) handleMsgWithdrawNative(origin string) error {
	return m.db.SaveHashes([]types.Hash{
		{
			Index: origin,
		},
	})
}

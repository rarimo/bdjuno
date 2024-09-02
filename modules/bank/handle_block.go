package bank

import (
	"fmt"
	juno "github.com/forbole/juno/v4/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	fmt.Println("bank....")

	err := m.saveTransfer(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating proposals: %s", err)
	}

	return nil
}

func (m *Module) saveTransfer(height int64) error {
	log.Debug().Str("module", m.Name()).Int64("height", height).
		Msg("saving transfer")

	events := m.keeper.GetEventsByBlockHeight(height)

	for _, event := range events {
		if err := m.db.SaveCoinTransfer(event.Amount, event.BlockHeight, event.Event, event.Address); err != nil {
			return err
		}
	}

	return nil
}

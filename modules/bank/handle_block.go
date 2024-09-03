package bank

import (
	juno "github.com/forbole/juno/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, txs []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	err := m.saveTransfer(txs, block.Block.Height)
	if err != nil {
		log.Error().Str("module", "bank").Err(err).Msg("error saving transfers to db")
		return errors.Wrap(err, "error saving transfers to db")
	}

	return nil
}

func (m *Module) saveTransfer(txs []*juno.Tx, height int64) error {
	log.Debug().Str("module", m.Name()).Int64("height", height).
		Msg("saving transfer")

	events, err := m.keeper.GetEventsFromTxs(txs)
	if err != nil {
		log.Error().Str("module", m.Name()).Err(err).Msg("error while getting events")
		return errors.Wrap(err, "error while getting events")
	}

	var commonError error
	for _, event := range events {
		if err = m.db.SaveCoinTransfer(event.Amount, event.BlockHeight, event.Event, event.Address, event.TxHash); err != nil {
			log.Error().Str("module", m.Name()).Err(err).Msg("error while saving transfer event")
			commonError = errors.Wrap(err, "error while saving transfer event")
		}
	}

	return commonError
}

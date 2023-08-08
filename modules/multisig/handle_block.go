package multisig

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/forbole/juno/v4/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	multisigtypes "gitlab.com/rarimo/rarimo-core/x/multisig/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	err := m.updateProposals(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating proposals: %s", err)
	}

	return nil
}

func (m *Module) updateProposals(height int64) error {
	log.Debug().Str("module", m.Name()).Int64("height", height).
		Msg("updating proposals")

	proposals, err := m.source.ProposalAll(height)
	if err != nil {
		return fmt.Errorf("error while getting proposals: %s", err)
	}

	for _, proposal := range proposals {
		if proposal.VotingEndBlock != uint64(height) {
			continue
		}

		if proposal.Status != multisigtypes.ProposalStatus_EXECUTED {
			continue
		}

		err = m.handleProposalExecution(height, proposal)
		if err != nil {
			return fmt.Errorf("error while handling proposal execution: %s", err)
		}
	}

	err = m.saveProposals(proposals)
	if err != nil {
		return fmt.Errorf("error while saving proposals: %s", err)
	}

	return nil
}

func (m *Module) handleProposalExecution(height int64, proposal multisigtypes.Proposal) error {
	addresses := make([]string, 0)

	for _, msg := range proposal.Messages {
		var stdMsg sdk.Msg
		err := m.cdc.UnpackAny(msg, &stdMsg)
		if err != nil {
			return fmt.Errorf("error while unpacking message: %s", err)
		}

		switch message := stdMsg.(type) {
		case *multisigtypes.MsgChangeGroup:
			err := m.handleMsgChangeGroup(height, message)
			if err != nil {
				return fmt.Errorf("error while handling multisend: %s", err)
			}
		case *banktypes.MsgSend:
			addresses = append(addresses, message.FromAddress, message.ToAddress)
		}
	}

	if len(addresses) == 0 {
		return nil
	}

	err := m.auth.RefreshAccounts(height, addresses)
	if err != nil {
		return fmt.Errorf("error while refreshing accounts: %s", err)
	}

	return nil
}

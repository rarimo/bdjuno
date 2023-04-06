package gov

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

func (m *Module) handleRarimoCoreProposal(height int64, proposal types.Proposal) error {
	if proposal.Status != types.StatusPassed {
		return nil
	}

	var content types.Content
	err := m.db.EncodingConfig.Marshaler.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling rarimo core proposal: %s", err)
	}

	switch content.(type) {
	case *rarimocoretypes.UnfreezeSignerPartyProposal,
		*rarimocoretypes.ReshareKeysProposal,
		*rarimocoretypes.ChangeThresholdProposal:
		return m.rarimocoreModule.UpdateParams(height)
	default:
		return nil
	}
}

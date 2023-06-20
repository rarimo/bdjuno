package gov

import (
	"fmt"
	govtypesv1beta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

func (m *Module) handleRarimoCoreProposal(height int64, proposal govtypesv1beta.Proposal) error {
	if proposal.Status != govtypesv1beta.StatusPassed {
		return nil
	}

	var content govtypesv1beta.Content
	err := m.db.EncodingConfig.Codec.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling rarimo core proposal: %s", err)
	}

	switch content.(type) {
	case *rarimocoretypes.UnfreezeSignerPartyProposal,
		*rarimocoretypes.ReshareKeysProposal,
		*rarimocoretypes.SlashProposal,
		*rarimocoretypes.DropPartiesProposal:
		return m.rarimocoreModule.UpdateParams(height)
	default:
		return nil
	}
}

package gov

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	oraclemanagertypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
)

func (m *Module) handleBridgeProposal(height int64, proposal govtypes.Proposal) error {
	if proposal.Status != govtypes.StatusPassed {
		return nil
	}

	var content govtypes.Content
	err := m.db.EncodingConfig.Codec.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling bridge proposal: %s", err)
	}

	switch content.(type) {
	case *oraclemanagertypes.ChangeParamsProposal:
		return m.bridgeModule.UpdateParams(height)
	default:
		return nil
	}
}

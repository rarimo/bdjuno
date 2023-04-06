package gov

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	oraclemanagertypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
)

func (m *Module) handleBridgeProposal(height int64, proposal types.Proposal) error {
	if proposal.Status != types.StatusPassed {
		return nil
	}

	var content types.Content
	err := m.db.EncodingConfig.Marshaler.UnpackAny(proposal.Content, &content)
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

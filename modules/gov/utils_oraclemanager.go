package gov

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	oraclemanagertypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
)

func (m *Module) handleOracleManagerProposal(height int64, proposal types.Proposal) error {
	if proposal.Status != types.StatusPassed {
		return nil
	}

	var content types.Content
	err := m.db.EncodingConfig.Marshaler.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling oracle manager proposal: %s", err)
	}

	switch cosmosMsg := content.(type) {
	case *oraclemanagertypes.OracleUnfreezeProposal:
		return m.oracleManagerModule.HandleOracle(height, cosmosMsg.Index.Chain, cosmosMsg.Index.Account)
	case *oraclemanagertypes.ChangeParamsProposal:
		return m.oracleManagerModule.UpdateParams(height)
	default:
		return nil
	}
}

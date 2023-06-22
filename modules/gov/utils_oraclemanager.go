package gov

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	oraclemanagertypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
)

func (m *Module) handleOracleManagerProposal(height int64, proposal govtypes.Proposal) error {
	if proposal.Status != govtypes.StatusPassed {
		return nil
	}

	var content govtypes.Content
	err := m.db.EncodingConfig.Codec.UnpackAny(proposal.Content, &content)
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

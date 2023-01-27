package gov

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

func (m *Module) handleTokenManagerProposal(height int64, rawProposal types.Proposal) error {
	if rawProposal.Status != types.StatusPassed {
		return nil
	}

	var content types.Content
	err := m.db.EncodingConfig.Marshaler.UnpackAny(rawProposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling tokenmanager proposal: %s", err)
	}

	switch proposal := content.(type) {
	case *tokenmanagertypes.SetNetworkProposal:
		return m.tokenmanagerModule.UpdateParams(height)
	case *tokenmanagertypes.UpdateTokenItemProposal:
		return m.tokenmanagerModule.UpdateItems(proposal.Item)
	case *tokenmanagertypes.RemoveTokenItemProposal:
		return m.tokenmanagerModule.RemoveItems(proposal.Index)
	case *tokenmanagertypes.CreateCollectionProposal:
		return m.tokenmanagerModule.CreateCollection(proposal.Index, proposal.Metadata, proposal.Data)
	case *tokenmanagertypes.UpdateCollectionDataProposal:
		return m.tokenmanagerModule.UpdateCollectionDatas(proposal.Data)
	case *tokenmanagertypes.AddCollectionDataProposal:
		return m.tokenmanagerModule.CreateCollectionDatas(height, proposal.Data)
	case *tokenmanagertypes.RemoveCollectionDataProposal:
		return m.tokenmanagerModule.RemoveCollectionDatas(height, proposal.Index)
	case *tokenmanagertypes.RemoveCollectionProposal:
		return m.tokenmanagerModule.RemoveCollection(proposal.Index)
	default:
		return nil
	}
}

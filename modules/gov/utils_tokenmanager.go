package gov

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

func (m *Module) handleTokenManagerProposal(height int64, proposal types.Proposal) error {
	if proposal.Status != types.StatusPassed {
		return nil
	}

	var content types.Content
	err := m.db.EncodingConfig.Marshaler.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling tokenmanager proposal: %s", err)
	}

	switch content.(type) {
	case *tokenmanagertypes.SetNetworkProposal:
		return m.tokenmanagerModule.UpdateParams(height)
	case *tokenmanagertypes.UpdateTokenItemProposal:
		return m.handleUpdateTokenItemProposal(content)
	case *tokenmanagertypes.RemoveTokenItemProposal:
		return m.handleRemoveTokenItemProposal(content)
	case *tokenmanagertypes.CreateCollectionProposal:
		return m.handleCreateCollectionProposal(content)
	case *tokenmanagertypes.UpdateCollectionDataProposal:
		return m.handleUpdateCollectionDataProposal(content)
	case *tokenmanagertypes.AddCollectionDataProposal:
		return m.handleAddCollectionDataProposal(content)
	case *tokenmanagertypes.RemoveCollectionDataProposal:
		return m.handleRemoveCollectionDataProposal(content)
	case *tokenmanagertypes.RemoveCollectionProposal:
		return m.handleRemoveCollectionProposal(content)
	default:
		return nil
	}
}

func (m *Module) handleUpdateTokenItemProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.UpdateTokenItemProposal)
	return m.tokenmanagerModule.UpdateItems(proposal.Item)
}

func (m *Module) handleRemoveTokenItemProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.RemoveTokenItemProposal)
	return m.tokenmanagerModule.RemoveItems(proposal.Index)
}

func (m *Module) handleCreateCollectionProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.CreateCollectionProposal)
	return m.tokenmanagerModule.CreateCollection(proposal.Index, proposal.Metadata, proposal.Data)
}

func (m *Module) handleUpdateCollectionDataProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.UpdateCollectionDataProposal)
	return m.tokenmanagerModule.UpdateCollectionDatas(proposal.Data)
}

func (m *Module) handleAddCollectionDataProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.AddCollectionDataProposal)
	return m.tokenmanagerModule.CreateCollectionDatas(proposal.Data)
}

func (m *Module) handleRemoveCollectionDataProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.RemoveCollectionDataProposal)
	return m.tokenmanagerModule.RemoveCollectionDatas(proposal.Index)
}

func (m *Module) handleRemoveCollectionProposal(content types.Content) error {
	proposal, _ := content.(*tokenmanagertypes.RemoveCollectionProposal)
	return m.tokenmanagerModule.RemoveCollection(proposal.Index)
}

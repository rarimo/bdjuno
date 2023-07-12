package gov

import (
	"fmt"
	govtypesv1beta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"math/big"
)

func (m *Module) handleTokenManagerProposal(height int64, rawProposal govtypesv1beta.Proposal) error {
	if rawProposal.Status != govtypesv1beta.StatusPassed {
		return nil
	}

	var content govtypesv1beta.Content
	err := m.db.EncodingConfig.Codec.UnpackAny(rawProposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling tokenmanager proposal: %s", err)
	}

	switch proposal := content.(type) {
	case *tokenmanagertypes.UpgradeContractProposal:
		return m.tokenmanagerModule.HandleUpdateContract(height, proposal.Details)
	case *tokenmanagertypes.AddNetworkProposal,
		*tokenmanagertypes.RemoveNetworkProposal:
		return m.tokenmanagerModule.UpdateParams(height)
	case *tokenmanagertypes.AddFeeTokenProposal:
		return m.handleAddOrUpdateFeeToken(height, proposal.Chain, proposal.Token)
	case *tokenmanagertypes.UpdateFeeTokenProposal:
		return m.handleAddOrUpdateFeeToken(height, proposal.Chain, proposal.Token)
	case *tokenmanagertypes.RemoveFeeTokenProposal:
		return m.handleRemoveFeeToken(height, proposal)
	case *tokenmanagertypes.WithdrawFeeProposal:
		return m.saveFeeManagementOp(height, proposal.Chain, proposal.Token.Contract, proposal.Token.Amount)
	case *tokenmanagertypes.UpdateTokenItemProposal:
		return m.tokenmanagerModule.UpdateItems(proposal.Item)
	case *tokenmanagertypes.RemoveTokenItemProposal:
		return m.tokenmanagerModule.RemoveItems(proposal.Index)
	case *tokenmanagertypes.CreateCollectionProposal:
		return m.tokenmanagerModule.CreateCollection(proposal.Index, &proposal.Metadata, proposal.Data, proposal.Item, proposal.OnChainItem)
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

func (m *Module) handleAddOrUpdateFeeToken(height int64, chain string, token tokenmanagertypes.FeeToken) error {
	err := m.tokenmanagerModule.UpdateParams(height)
	if err != nil {
		return fmt.Errorf("error while updating tokenmanager proposal: %s", err)
	}

	return m.saveFeeManagementOp(height, chain, token.Contract, token.Amount)
}

func (m *Module) handleRemoveFeeToken(height int64, proposal *tokenmanagertypes.RemoveFeeTokenProposal) error {
	err := m.tokenmanagerModule.UpdateParams(height)
	if err != nil {
		return fmt.Errorf("error while updating tokenmanager proposal: %s", err)
	}

	token, err := m.tokenmanagerModule.GetFeeToken(height, proposal.Chain, proposal.Contract)
	if err != nil {
		return fmt.Errorf("error while getting fee token to remove: %s", err)
	}
	return m.saveFeeManagementOp(height, proposal.Chain, token.Contract, token.Amount)
}

func (m *Module) saveFeeManagementOp(height int64, chain, contract string, amount string) error {
	index := hexutil.Encode(crypto.Keccak256(big.NewInt(height).Bytes(), []byte(chain), []byte(contract), []byte(amount)))
	err := m.rarimocoreModule.SaveOperationByIndex(height, index)
	if err != nil {
		return fmt.Errorf("error while saving operation: %s", err)
	}
	return nil
}

package rarimocore

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"gitlab.com/rarimo/bdjuno/types"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

func (m *Module) SaveOperationByIndex(height int64, index string) error {
	operation, err := m.source.Operation(height, index)
	if err != nil {
		return fmt.Errorf("failed to get operation by index: %s", err)
	}

	err = m.saveOperations([]rarimocoretypes.Operation{operation})
	if err != nil {
		return fmt.Errorf("failed to save operation: %s", err)
	}
	return nil
}

func (m *Module) saveOperations(slice []rarimocoretypes.Operation) error {
	operations := coreOperationsToInternal(slice)

	// Save the operations
	err := m.db.SaveOperations(operations)
	if err != nil {
		return err
	}

	transfers, changeParties, contractUpgrades, feeTokenManagements, identityTransfers, err := getOperationDetails(slice)
	if err != nil {
		return nil
	}

	// Save the transfers
	err = m.db.SaveTransfers(transfers)
	if err != nil {
		return err
	}

	// Save the change parties
	err = m.db.SaveChangeParties(changeParties)
	if err != nil {
		return err
	}

	// Save the contract upgrades
	err = m.db.SaveContractUpgrades(contractUpgrades)
	if err != nil {
		return err
	}

	// Save the fee token managements
	err = m.db.SaveFeeTokenManagements(feeTokenManagements)
	if err != nil {
		return err
	}

	// Save the identity transfers
	err = m.db.SaveIdentityDefaultTransfers(identityTransfers)
	if err != nil {
		return err
	}

	return nil

}

func (m *Module) saveConfirmations(slice []rarimocoretypes.Confirmation) error {
	confirmations := make([]types.Confirmation, len(slice))
	for index, confirmation := range slice {
		confirmations[index] = types.NewConfirmation(confirmation)
	}
	return m.db.SaveConfirmations(confirmations)
}

func (m *Module) saveVotes(slice []rarimocoretypes.Vote) error {
	votes := make([]types.RarimoCoreVote, len(slice))
	for index, vote := range slice {
		votes[index] = types.RarimoCoreVoteFromCore(vote)
	}
	return m.db.SaveRarimoCoreVotes(votes)
}

func (m *Module) updateOperations(slice []rarimocoretypes.Operation) error {
	operations := coreOperationsToInternal(slice)

	// Update the operations
	for _, operation := range operations {
		err := m.db.UpdateOperation(operation)
		if err != nil {
			return err
		}
	}

	_, changeParties, _, _, _, err := getOperationDetails(slice)
	if err != nil {
		return err
	}

	// Update the change parties
	for _, changeParty := range changeParties {
		err = m.db.UpdateChangeParties(changeParty)
		if err != nil {
			return err
		}
	}

	return nil

}

func coreOperationsToInternal(slice []rarimocoretypes.Operation) []types.Operation {
	operations := make([]types.Operation, len(slice))
	for i, operation := range slice {
		operations[i] = types.OperationFromCore(operation)
	}

	return operations
}

func getOperationDetails(slice []rarimocoretypes.Operation) ([]types.Transfer, []types.ChangeParties, []types.ContractUpgrade, []types.FeeTokenManagement, []types.IdentityDefaultTransfer, error) {
	transfers := make([]types.Transfer, 0)
	changeParties := make([]types.ChangeParties, 0)
	contractUpdates := make([]types.ContractUpgrade, 0)
	feeTokenManagements := make([]types.FeeTokenManagement, 0)
	identityTransfers := make([]types.IdentityDefaultTransfer, 0)

	for _, operation := range slice {
		switch operation.OperationType {
		case rarimocoretypes.OpType_TRANSFER:
			transfer, err := getTransfer(operation)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			transfers = append(transfers, transfer)
		case rarimocoretypes.OpType_CHANGE_PARTIES:
			changeParty, err := getChangeParties(operation)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			changeParties = append(changeParties, changeParty)
		case rarimocoretypes.OpType_CONTRACT_UPGRADE:
			contractUpdate, err := getContractUpdate(operation)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			contractUpdates = append(contractUpdates, contractUpdate)
		case rarimocoretypes.OpType_FEE_TOKEN_MANAGEMENT:
			feeTokenManagement, err := getFeeTokenManagement(operation)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			feeTokenManagements = append(feeTokenManagements, feeTokenManagement)
		case rarimocoretypes.OpType_IDENTITY_DEFAULT_TRANSFER:
			identityDefaultTransfer, err := getIdentityDefaultTransfer(operation)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			identityTransfers = append(identityTransfers, identityDefaultTransfer)
		default:
			log.Warn().Str("module", "rarimocore").
				Str("operation_type", string(operation.OperationType)).
				Msg("unknown operation type")
			return nil, nil, nil, nil, nil, errors.Wrap(errors.ErrInvalidType, "invalid operation type")
		}
	}

	return transfers, changeParties, contractUpdates, feeTokenManagements, identityTransfers, nil
}

func getTransfer(operation rarimocoretypes.Operation) (types.Transfer, error) {
	transfer := new(rarimocoretypes.Transfer)
	err := proto.Unmarshal(operation.Details.Value, transfer)
	if err != nil {
		return types.Transfer{}, err
	}

	return types.NewTransfer(operation.Index, *transfer), nil
}

func getChangeParties(operation rarimocoretypes.Operation) (types.ChangeParties, error) {
	changeParties := new(rarimocoretypes.ChangeParties)
	err := proto.Unmarshal(operation.Details.Value, changeParties)
	if err != nil {
		return types.ChangeParties{}, err
	}

	return types.NewChangeParties(operation.Index, *changeParties), nil
}

func getContractUpdate(operation rarimocoretypes.Operation) (types.ContractUpgrade, error) {
	contractUpdate := new(rarimocoretypes.ContractUpgrade)
	err := proto.Unmarshal(operation.Details.Value, contractUpdate)
	if err != nil {
		return types.ContractUpgrade{}, err
	}

	return types.NewContractUpdate(operation.Index, *contractUpdate), nil
}

func getFeeTokenManagement(operation rarimocoretypes.Operation) (types.FeeTokenManagement, error) {
	feeTokenManagement := new(rarimocoretypes.FeeTokenManagement)
	err := proto.Unmarshal(operation.Details.Value, feeTokenManagement)
	if err != nil {
		return types.FeeTokenManagement{}, err
	}

	return types.NewFeeTokenManagement(operation.Index, *feeTokenManagement), nil
}

func getIdentityDefaultTransfer(operation rarimocoretypes.Operation) (types.IdentityDefaultTransfer, error) {
	identityTransfer := new(rarimocoretypes.IdentityDefaultTransfer)
	err := proto.Unmarshal(operation.Details.Value, identityTransfer)
	if err != nil {
		return types.IdentityDefaultTransfer{}, err
	}

	return types.NewIdentityDefaultTransfer(operation.Index, *identityTransfer), nil
}

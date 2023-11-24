package rarimocore

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rarimo/bdjuno/types"
	rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/rs/zerolog/log"
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

	transfers := make([]types.Transfer, 0)
	changeParties := make([]types.ChangeParties, 0)
	contractUpdates := make([]types.ContractUpgrade, 0)
	feeTokenManagements := make([]types.FeeTokenManagement, 0)
	identityDefaultTransfers := make([]types.IdentityDefaultTransfer, 0)
	identityGISTTransfers := make([]types.IdentityGISTTransfer, 0)
	identityStateTransfers := make([]types.IdentityStateTransfer, 0)

	for _, operation := range slice {
		switch operation.OperationType {
		case rarimocoretypes.OpType_TRANSFER:
			transfer, err := getTransfer(operation)
			if err != nil {
				return fmt.Errorf("error while extracting transfer from operation: %s", err)
			}
			transfers = append(transfers, transfer)
		case rarimocoretypes.OpType_CHANGE_PARTIES:
			changeParty, err := getChangeParties(operation)
			if err != nil {
				return fmt.Errorf("error while extracting change parties from operation: %s", err)
			}
			changeParties = append(changeParties, changeParty)
		case rarimocoretypes.OpType_CONTRACT_UPGRADE:
			contractUpdate, err := getContractUpdate(operation)
			if err != nil {
				return fmt.Errorf("error while extracting contract updates from operation: %s", err)
			}
			contractUpdates = append(contractUpdates, contractUpdate)
		case rarimocoretypes.OpType_FEE_TOKEN_MANAGEMENT:
			feeTokenManagement, err := getFeeTokenManagement(operation)
			if err != nil {
				return fmt.Errorf("error while extracting fee token management from operation: %s", err)
			}
			feeTokenManagements = append(feeTokenManagements, feeTokenManagement)
		case rarimocoretypes.OpType_IDENTITY_DEFAULT_TRANSFER:
			identityDefaultTransfer, err := getIdentityDefaultTransfer(operation)
			if err != nil {
				return fmt.Errorf("error while extracting identity default transfer from operation: %s", err)
			}
			identityDefaultTransfers = append(identityDefaultTransfers, identityDefaultTransfer)
		case rarimocoretypes.OpType_IDENTITY_GIST_TRANSFER:
			identityGISTTransfer, err := getIdentityGISTTransfer(operation)
			if err != nil {
				return fmt.Errorf("error while extracting identity gist transfer from operation: %s", err)
			}
			identityGISTTransfers = append(identityGISTTransfers, identityGISTTransfer)
		case rarimocoretypes.OpType_IDENTITY_STATE_TRANSFER:
			identityStateTransfer, err := getIdentityStateTransfer(operation)
			if err != nil {
				return fmt.Errorf("error while extracting identity state transfer from operation: %s", err)
			}
			identityStateTransfers = append(identityStateTransfers, identityStateTransfer)
		default:
			log.Warn().Str("module", "rarimocore").
				Str("operation_type", string(operation.OperationType)).
				Msg("unknown operation type")
			return fmt.Errorf("error while handling operation content %s: ", errors.ErrInvalidType)
		}
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
	err = m.db.SaveContractUpgrades(contractUpdates)
	if err != nil {
		return err
	}

	// Save the fee token managements
	err = m.db.SaveFeeTokenManagements(feeTokenManagements)
	if err != nil {
		return err
	}

	// Save the identity default transfers
	err = m.db.SaveIdentityDefaultTransfers(identityDefaultTransfers)
	if err != nil {
		return err
	}

	// Save the identity gist transfers
	err = m.db.SaveIdentityGISTTransfers(identityGISTTransfers)
	if err != nil {
		return err
	}

	// Save the identity state transfers
	err = m.db.SaveIdentityStateTransfers(identityStateTransfers)
	if err != nil {
		return err
	}

	return nil

}

func (m *Module) saveConfirmations(slice []rarimocoretypes.Confirmation, height int64, tx *string) error {
	confirmations := make([]types.Confirmation, len(slice))
	for index, confirmation := range slice {
		confirmations[index] = types.NewConfirmation(confirmation, height, tx)
	}
	return m.db.SaveConfirmations(confirmations)
}

func (m *Module) saveVotes(slice []rarimocoretypes.Vote, height int64, tx *string) error {
	votes := make([]types.RarimoCoreVote, len(slice))
	for index, vote := range slice {
		votes[index] = types.RarimoCoreVoteFromCore(vote, height, tx)
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

	for _, operation := range slice {
		switch operation.OperationType {
		case rarimocoretypes.OpType_CHANGE_PARTIES:
			changeParty, err := getChangeParties(operation)
			if err != nil {
				return err
			}
			err = m.db.UpdateChangeParties(changeParty)
			if err != nil {
				return err
			}
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

func getIdentityGISTTransfer(operation rarimocoretypes.Operation) (types.IdentityGISTTransfer, error) {
	transfer := new(rarimocoretypes.IdentityGISTTransfer)
	err := proto.Unmarshal(operation.Details.Value, transfer)
	if err != nil {
		return types.IdentityGISTTransfer{}, err
	}

	return types.NewIdentityGISTTransfer(operation.Index, *transfer), nil
}

func getIdentityStateTransfer(operation rarimocoretypes.Operation) (types.IdentityStateTransfer, error) {
	transfer := new(rarimocoretypes.IdentityStateTransfer)
	err := proto.Unmarshal(operation.Details.Value, transfer)
	if err != nil {
		return types.IdentityStateTransfer{}, err
	}

	return types.NewIdentityStateTransfer(operation.Index, *transfer), nil
}

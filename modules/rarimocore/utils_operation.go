package rarimocore

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"gitlab.com/rarimo/bdjuno/types"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

func (m *Module) saveOperations(slice []rarimocoretypes.Operation) error {
	operations := coreOperationsToInternal(slice)

	return m.db.Transaction(func() error {
		// Save the operations
		err := m.db.SaveOperations(operations)
		if err != nil {
			return err
		}

		transfers, changeParties, err := getOperationDetails(slice)
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

		return nil
	})
}

func (m *Module) saveConfirmations(slice []rarimocoretypes.Confirmation) error {
	confirmations := make([]types.Confirmation, len(slice))
	for index, confirmation := range slice {
		confirmations[index] = types.NewConfirmation(confirmation)
	}
	return m.db.SaveConfirmations(confirmations)
}

func (m *Module) updateOperations(slice []rarimocoretypes.Operation) error {
	operations := coreOperationsToInternal(slice)

	return m.db.Transaction(func() error {
		// Update the operations
		for _, operation := range operations {
			err := m.db.UpdateOperation(operation)
			if err != nil {
				return err
			}
		}

		_, changeParties, err := getOperationDetails(slice)
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
	})
}

func coreOperationsToInternal(slice []rarimocoretypes.Operation) []types.Operation {
	operations := make([]types.Operation, len(slice))
	for i, operation := range slice {
		operations[i] = types.OperationFromCore(operation)
	}

	return operations
}

func getOperationDetails(slice []rarimocoretypes.Operation) ([]types.Transfer, []types.ChangeParties, error) {
	transfers := make([]types.Transfer, 0)
	changeParties := make([]types.ChangeParties, 0)

	for _, operation := range slice {
		switch operation.OperationType {
		case rarimocoretypes.OpType_TRANSFER:
			transfer, err := getTransfer(operation)
			if err != nil {
				return nil, nil, err
			}
			transfers = append(transfers, transfer)
		case rarimocoretypes.OpType_CHANGE_PARTIES:
			changeParty, err := getChangeParties(operation)
			if err != nil {
				return nil, nil, err
			}
			changeParties = append(changeParties, changeParty)
		default:
			log.Warn().Str("module", "rarimocore").
				Str("operation_type", string(operation.OperationType)).
				Msg("unknown operation type")
			return nil, nil, errors.Wrap(errors.ErrInvalidType, "invalid operation type")
		}
	}

	return transfers, changeParties, nil
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

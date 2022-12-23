package rarimocore

import (
	"encoding/json"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"gitlab.com/rarimo/bdjuno/types"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "rarimocore").Msg("parsing genesis")

	// Read the genesis state
	var genState rarimocoretypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[rarimocoretypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading rarimocore genesis data: %s", err)
	}

	// Save the operations
	err = m.saveOperations(genState.OperationList)
	if err != nil {
		return fmt.Errorf("error while storing genesis rarimocore operations: %s", err)
	}

	// Save the confirmations
	err = m.saveConfirmations(genState.ConfirmationList)
	if err != nil {
		return fmt.Errorf("error while storing genesis rarimocore confirmations: %s", err)
	}

	// Save the params
	err = m.saveParams(genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis rarimocore params: %s", err)
	}

	// Save the parties
	err = m.saveParties(genState.Params.Parties)
	if err != nil {
		return fmt.Errorf("error while storing genesis rarimocore parties: %s", err)
	}

	return nil
}

func (m *Module) saveParties(slice []*rarimocoretypes.Party) error {
	parties := make([]types.Party, len(slice))
	for index, party := range slice {
		// We don't have operation index for initial change parties from genesis
		parties[index] = types.NewParty(*party)
	}
	return m.db.SaveParties(parties)
}

func (m *Module) saveParams(params rarimocoretypes.Params, height int64) (err error) {
	p := types.NewRarimoCoreParams(params, height)
	err = m.db.SaveRarimoCoreParams(p)
	if err != nil {
		return err
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

	transfers, changeParties, err := m.getOperationDetails(slice)
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
}

func (m *Module) saveConfirmations(slice []rarimocoretypes.Confirmation) error {
	confirmations := make([]types.Confirmation, len(slice))
	for index, confirmation := range slice {
		confirmations[index] = types.NewConfirmation(confirmation)
	}
	return m.db.SaveConfirmations(confirmations)
}

func (m *Module) getOperationDetails(slice []rarimocoretypes.Operation) ([]types.Transfer, []types.ChangeParties, error) {
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
			return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid operation type")
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

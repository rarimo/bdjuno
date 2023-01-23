package rarimocore

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState rarimocoretypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[rarimocoretypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading rarimocore genesis data: %s", err)
	}

	return m.db.Transaction(func() error {
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
	})
}
package oraclemanager

import (
	"encoding/json"
	"fmt"
	"github.com/rarimo/bdjuno/types"
	oraclemanagertypes "github.com/rarimo/rarimo-core/x/oraclemanager/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState oraclemanagertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[oraclemanagertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading oraclemanager genesis data: %s", err)
	}
	// Save the params
	err = m.saveParams(genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis oraclemanager params: %s", err)
	}

	oracles := make([]types.Oracle, len(genState.Oracles))
	for i, oracle := range genState.Oracles {
		oracles[i] = types.OracleFromCore(oracle)
	}

	err = m.db.SaveOracles(oracles)
	if err != nil {
		return fmt.Errorf("error while storing genesis oracles: %s", err)
	}

	return nil

}

package bridge

import (
	"encoding/json"
	"fmt"
	"github.com/rarimo/bdjuno/types"
	bridgetypes "github.com/rarimo/rarimo-core/x/bridge/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState bridgetypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[bridgetypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading bridge genesis data: %s", err)
	}
	// Save the params
	err = m.saveParams(genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis bridge params: %s", err)
	}

	hashes := make([]types.Hash, len(genState.HashList))
	for i, hash := range genState.HashList {
		hashes[i] = types.HashFromCore(hash)
	}

	err = m.db.SaveHashes(hashes)
	if err != nil {
		return fmt.Errorf("error while storing genesis bridge hashes: %s", err)
	}

	return nil

}

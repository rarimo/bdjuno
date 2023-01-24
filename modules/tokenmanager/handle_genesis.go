package tokenmanager

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState tokenmanagertypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[tokenmanagertypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading tokenmanager genesis data: %s", err)
	}

	return m.db.Transaction(func() error {
		// Save the collection datas
		err = m.saveCollectionDatas(genState.Datas)
		if err != nil {
			return fmt.Errorf("error while storing genesis tokenmanager collection datas: %s", err)
		}

		// Save the collections
		err = m.saveCollections(genState.Collections)
		if err != nil {
			return fmt.Errorf("error while storing genesis tokenmanager collections: %s", err)
		}

		// Save the items
		err = m.saveItems(genState.Items)
		if err != nil {
			return fmt.Errorf("error while storing genesis tokenmanager items: %s", err)
		}

		// Save the params
		err = m.saveParams(genState.Params, doc.InitialHeight)
		if err != nil {
			return fmt.Errorf("error while storing genesis tokenmanager params: %s", err)
		}

		return nil
	})
}

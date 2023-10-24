package multisig

import (
	"encoding/json"
	"fmt"
	"github.com/rarimo/bdjuno/types"
	multisigtypes "github.com/rarimo/rarimo-core/x/multisig/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", m.Name()).Msg("parsing genesis")

	// Read the genesis state
	var genState multisigtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[multisigtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading multisig genesis data: %s", err)
	}

	err = m.db.SaveMultisigParams(types.MultisigParamsFromCore(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis multisig params: %s", err)
	}

	err = m.saveGroups(genState.GroupList)
	if err != nil {
		return fmt.Errorf("error while storing genesis multisig groups: %s", err)
	}

	err = m.saveProposals(genState.ProposalList)
	if err != nil {
		return fmt.Errorf("error while storing genesis multisig proposals: %s", err)
	}

	err = m.saveVotes(genState.VoteList)
	if err != nil {
		return fmt.Errorf("error while storing genesis multisig proposal votes: %s", err)
	}

	return nil

}

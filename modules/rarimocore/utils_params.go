package rarimocore

import (
	"fmt"
	"github.com/rarimo/bdjuno/types"
	rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"
	"github.com/rs/zerolog/log"
)

func (m *Module) UpdateParams(height int64) error {
	log.Debug().Str("module",
		m.Name()).Int64("height", height).
		Msg("updating params")

	params, err := m.source.Params(height)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	err = m.saveParams(params, height)
	if err != nil {
		return fmt.Errorf("error while storing params during update rarimocore params: %s", err)
	}

	err = m.saveParties(params.Parties)
	if err != nil {
		return fmt.Errorf("error while storing parties during update rarimocore params: %s", err)
	}

	return nil
}

func (m *Module) saveParties(slice []*rarimocoretypes.Party) error {
	parties := make([]types.Party, len(slice))
	for index, party := range slice {
		parties[index] = types.NewParty(*party)
	}
	return m.db.SaveParties(parties)
}

func (m *Module) saveParams(params rarimocoretypes.Params, height int64) (err error) {
	return m.db.SaveRarimoCoreParams(types.NewRarimoCoreParams(params, height))
}

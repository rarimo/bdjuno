package tokenmanager

import (
	"fmt"
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

	err = m.saveNetworks(params)
	if err != nil {
		return fmt.Errorf("error while storing params during update tokenmanager networks: %s", err)
	}

	return nil
}

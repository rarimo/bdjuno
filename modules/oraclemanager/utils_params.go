package oraclemanager

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab.com/rarimo/bdjuno/types"
	oraclemanagertypes "gitlab.com/rarimo/rarimo-core/x/oraclemanager/types"
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
		return fmt.Errorf("error while storing params during update oraclemanager params: %s", err)
	}

	return nil
}

func (m *Module) saveParams(params oraclemanagertypes.Params, height int64) (err error) {
	return m.db.SaveOracleManagerParams(types.OracleManagerParamsFromCore(params, height))
}

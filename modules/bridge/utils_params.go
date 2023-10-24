package bridge

import (
	"fmt"
	"github.com/rarimo/bdjuno/types"
	bridgetypes "github.com/rarimo/rarimo-core/x/bridge/types"
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
		return fmt.Errorf("error while storing params during update bridge params: %s", err)
	}

	return nil
}

func (m *Module) saveParams(params bridgetypes.Params, height int64) (err error) {
	return m.db.SaveBridgeParams(types.BridgeParamsFromCore(params, height))
}

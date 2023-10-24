package source

import bridgetypes "github.com/rarimo/rarimo-core/x/bridge/types"

type Source interface {
	Params(height int64) (bridgetypes.Params, error)
}

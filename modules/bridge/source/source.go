package source

import bridgetypes "gitlab.com/rarimo/rarimo-core/x/bridge/types"

type Source interface {
	Params(height int64) (bridgetypes.Params, error)
}

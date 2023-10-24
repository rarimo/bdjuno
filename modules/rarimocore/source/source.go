package source

import rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"

type Source interface {
	Params(height int64) (rarimocoretypes.Params, error)
	Operation(height int64, index string) (rarimocoretypes.Operation, error)
	Confirmation(height int64, root string) (rarimocoretypes.Confirmation, error)
	ViolationReport(height int64, sessionId, offender, sender string, violationType rarimocoretypes.ViolationType) (rarimocoretypes.ViolationReport, error)
}

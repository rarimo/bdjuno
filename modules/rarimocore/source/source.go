package source

import rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"

type Source interface {
	Params(height int64) (rarimocoretypes.Params, error)
	Operation(height int64, index string) (rarimocoretypes.Operation, error)
	OperationAll(height int64) ([]rarimocoretypes.Operation, error)
	Confirmation(height int64, root string) (rarimocoretypes.Confirmation, error)
	ConfirmationAll(height int64) ([]rarimocoretypes.Confirmation, error)
}

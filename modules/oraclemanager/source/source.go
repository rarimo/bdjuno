package source

import oraclemanagertypes "github.com/rarimo/rarimo-core/x/oraclemanager/types"

type Source interface {
	Params(height int64) (oraclemanagertypes.Params, error)
	Oracle(height int64, chain, account string) (oraclemanagertypes.Oracle, error)
}

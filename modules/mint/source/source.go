package source

import (
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
}

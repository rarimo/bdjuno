package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"gitlab.com/rarimo/bdjuno/types"
)

type SlashingModule interface {
	GetSigningInfo(height int64, consAddr sdk.ConsAddress) (types.ValidatorSigningInfo, error)
}

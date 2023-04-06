package types

import (
	bridgetypes "gitlab.com/rarimo/rarimo-core/x/bridge/types"
)

// BridgeParams contains the data of the x/bridge module params instance
type BridgeParams struct {
	WithdrawDenom string `json:"withdraw_denom,omitempty" yaml:"withdraw_denom,omitempty"`
	Height        int64  `json:"height,omitempty" yaml:"height,omitempty"`
}

// BridgeParamsFromCore allows to build a new BridgeParams instance from an oraclemanagertypes.Params instance
func BridgeParamsFromCore(p bridgetypes.Params, height int64) *BridgeParams {
	return &BridgeParams{
		WithdrawDenom: p.WithdrawDenom,
		Height:        height,
	}
}

// Hash represents a single hash instance
type Hash struct {
	Index string `json:"index,omitempty" yaml:"index,omitempty"`
}

// HashFromCore allows to build a new Hash instance from an bridgetypes.Hash instance
func HashFromCore(hash bridgetypes.Hash) Hash {
	return Hash{
		Index: hash.Index,
	}
}

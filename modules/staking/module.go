package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"gitlab.com/rarimo/bdjuno/database"
	stakingsource "gitlab.com/rarimo/bdjuno/modules/staking/source"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/staking module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source stakingsource.Source
}

// NewModule returns a new Module instance
func NewModule(
	source stakingsource.Source, cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "staking"
}

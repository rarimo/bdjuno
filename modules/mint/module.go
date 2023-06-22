package mint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"

	"gitlab.com/rarimo/bdjuno/database"
	mintsource "gitlab.com/rarimo/bdjuno/modules/mint/source"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represent database/mint module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source mintsource.Source
}

// NewModule returns a new Module instance
func NewModule(source mintsource.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "mint"
}

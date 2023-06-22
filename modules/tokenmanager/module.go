package tokenmanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	tokenmanager "gitlab.com/rarimo/bdjuno/modules/tokenmanager/source"

	"gitlab.com/rarimo/bdjuno/database"

	"github.com/forbole/juno/v4/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
)

// Module represents the x/tokenmanager module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source tokenmanager.Source
}

// NewModule builds a new Module instance
func NewModule(source tokenmanager.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		source: source,
		cdc:    cdc,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "tokenmanager"
}

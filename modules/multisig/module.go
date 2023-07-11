package multisig

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"gitlab.com/rarimo/bdjuno/database"
	multisig "gitlab.com/rarimo/bdjuno/modules/multisig/source"

	"github.com/forbole/juno/v4/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
	_ modules.BlockModule   = &Module{}
)

// Module represents the x/multisig module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source multisig.Source
}

// NewModule builds a new Module instance
func NewModule(source multisig.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		source: source,
		cdc:    cdc,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "multisig"
}

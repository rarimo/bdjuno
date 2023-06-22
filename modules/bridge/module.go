package bridge

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"gitlab.com/rarimo/bdjuno/database"
	bridge "gitlab.com/rarimo/bdjuno/modules/bridge/source"

	"github.com/forbole/juno/v4/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/bridge module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source bridge.Source
}

// NewModule builds a new Module instance
func NewModule(source bridge.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		source: source,
		cdc:    cdc,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "bridge"
}

package rarimocore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	rarimocore "gitlab.com/rarimo/bdjuno/modules/rarimocore/source"
	tokenmanager "gitlab.com/rarimo/bdjuno/modules/tokenmanager/source"

	"gitlab.com/rarimo/bdjuno/database"

	"github.com/forbole/juno/v4/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/rarimocore module
type Module struct {
	cdc                codec.Codec
	db                 *database.Db
	source             rarimocore.Source
	tokenmanagerSource tokenmanager.Source
}

// NewModule builds a new Module instance
func NewModule(source rarimocore.Source, tokenmanagerSource tokenmanager.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		source:             source,
		tokenmanagerSource: tokenmanagerSource,
		cdc:                cdc,
		db:                 db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "rarimocore"
}

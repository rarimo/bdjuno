package oraclemanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"gitlab.com/rarimo/bdjuno/database"
	oraclemanager "gitlab.com/rarimo/bdjuno/modules/oraclemanager/source"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/oraclemanager module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source oraclemanager.Source
}

// NewModule builds a new Module instance
func NewModule(source oraclemanager.Source, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		source: source,
		cdc:    cdc,
		db:     db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "oraclemanager"
}

package oraclemanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"gitlab.com/rarimo/bdjuno/database"
	oraclemanager "gitlab.com/rarimo/bdjuno/modules/oraclemanager/source"
	"gitlab.com/rarimo/bdjuno/modules/rarimocore"

	"github.com/forbole/juno/v4/modules"
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
	rc     *rarimocore.Module
}

// NewModule builds a new Module instance
func NewModule(source oraclemanager.Source, rc *rarimocore.Module, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		source: source,
		cdc:    cdc,
		db:     db,
		rc:     rc,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "oraclemanager"
}

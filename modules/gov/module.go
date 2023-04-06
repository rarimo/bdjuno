package gov

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"gitlab.com/rarimo/bdjuno/database"

	govsource "gitlab.com/rarimo/bdjuno/modules/gov/source"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.BlockModule   = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represent x/gov module
type Module struct {
	cdc                 codec.Codec
	db                  *database.Db
	source              govsource.Source
	authModule          AuthModule
	distrModule         DistrModule
	mintModule          MintModule
	slashingModule      SlashingModule
	stakingModule       StakingModule
	rarimocoreModule    RarimoCoreModule
	tokenmanagerModule  TokenManagerModule
	oracleManagerModule OracleManagerModule
}

// NewModule returns a new Module instance
func NewModule(
	source govsource.Source,
	authModule AuthModule,
	distrModule DistrModule,
	mintModule MintModule,
	slashingModule SlashingModule,
	stakingModule StakingModule,
	rarimocoreModule RarimoCoreModule,
	tokenmanagerModule TokenManagerModule,
	oracleManager OracleManagerModule,
	cdc codec.Codec,
	db *database.Db,
) *Module {
	return &Module{
		cdc:                 cdc,
		source:              source,
		authModule:          authModule,
		distrModule:         distrModule,
		mintModule:          mintModule,
		slashingModule:      slashingModule,
		stakingModule:       stakingModule,
		rarimocoreModule:    rarimocoreModule,
		tokenmanagerModule:  tokenmanagerModule,
		oracleManagerModule: oracleManager,
		db:                  db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "gov"
}

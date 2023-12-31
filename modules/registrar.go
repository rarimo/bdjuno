package modules

import (
	"github.com/rarimo/bdjuno/modules/actions"
	"github.com/rarimo/bdjuno/modules/bridge"
	"github.com/rarimo/bdjuno/modules/multisig"
	"github.com/rarimo/bdjuno/modules/oraclemanager"
	"github.com/rarimo/bdjuno/modules/rarimocore"
	"github.com/rarimo/bdjuno/modules/tokenmanager"
	"github.com/rarimo/bdjuno/modules/types"

	"github.com/forbole/juno/v4/modules/pruning"
	"github.com/forbole/juno/v4/modules/telemetry"

	"github.com/rarimo/bdjuno/modules/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jmodules "github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/modules/messages"
	"github.com/forbole/juno/v4/modules/registrar"

	"github.com/rarimo/bdjuno/utils"

	"github.com/rarimo/bdjuno/database"
	"github.com/rarimo/bdjuno/modules/auth"
	"github.com/rarimo/bdjuno/modules/bank"
	"github.com/rarimo/bdjuno/modules/consensus"
	"github.com/rarimo/bdjuno/modules/distribution"
	"github.com/rarimo/bdjuno/modules/feegrant"

	"github.com/rarimo/bdjuno/modules/gov"
	"github.com/rarimo/bdjuno/modules/mint"
	"github.com/rarimo/bdjuno/modules/modules"
	"github.com/rarimo/bdjuno/modules/pricefeed"
	"github.com/rarimo/bdjuno/modules/staking"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(cdc codec.Codec, msg sdk.Msg) ([]string, error) {
		addresses, err := parser(cdc, msg)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := ctx.EncodingConfig.Codec
	db := database.Cast(ctx.Database)

	sources, err := types.BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	actionsModule := actions.NewModule(ctx.JunoConfig, ctx.EncodingConfig)
	authModule := auth.NewModule(r.parser, cdc, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	consensusModule := consensus.NewModule(db)
	distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
	feegrantModule := feegrant.NewModule(cdc, db)
	mintModule := mint.NewModule(sources.MintSource, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, cdc, db)
	rarimocoreModule := rarimocore.NewModule(sources.RarimoCoreSource, sources.TokenManagerSource, cdc, db)
	tokenmanagerModule := tokenmanager.NewModule(sources.TokenManagerSource, cdc, db)
	oraclemanagerModule := oraclemanager.NewModule(sources.OracleManagerSource, rarimocoreModule, cdc, db)
	bridgeModule := bridge.NewModule(sources.BridgeSource, cdc, db)
	govModule := gov.NewModule(
		sources.GovSource,
		authModule,
		distrModule,
		mintModule,
		slashingModule,
		stakingModule,
		rarimocoreModule,
		tokenmanagerModule,
		oraclemanagerModule,
		bridgeModule,
		cdc,
		db,
	)

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		actionsModule,
		authModule,
		bankModule,
		consensusModule,
		distrModule,
		feegrantModule,
		govModule,
		mintModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		stakingModule,
		rarimocoreModule,
		tokenmanagerModule,
		oraclemanagerModule,
		bridgeModule,
		multisig.NewModule(sources.MultisigSource, cdc, db, authModule),
	}
}

package modules

import (
	"gitlab.com/rarimo/bdjuno/modules/actions"
	"gitlab.com/rarimo/bdjuno/modules/rarimocore"
	"gitlab.com/rarimo/bdjuno/modules/tokenmanager"
	"gitlab.com/rarimo/bdjuno/modules/types"

	"github.com/forbole/juno/v3/modules/pruning"
	"github.com/forbole/juno/v3/modules/telemetry"

	"gitlab.com/rarimo/bdjuno/modules/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jmodules "github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/modules/messages"
	"github.com/forbole/juno/v3/modules/registrar"

	"gitlab.com/rarimo/bdjuno/utils"

	"gitlab.com/rarimo/bdjuno/database"
	"gitlab.com/rarimo/bdjuno/modules/auth"
	"gitlab.com/rarimo/bdjuno/modules/bank"
	"gitlab.com/rarimo/bdjuno/modules/consensus"
	"gitlab.com/rarimo/bdjuno/modules/distribution"
	"gitlab.com/rarimo/bdjuno/modules/feegrant"

	"gitlab.com/rarimo/bdjuno/modules/gov"
	"gitlab.com/rarimo/bdjuno/modules/mint"
	"gitlab.com/rarimo/bdjuno/modules/modules"
	"gitlab.com/rarimo/bdjuno/modules/pricefeed"
	"gitlab.com/rarimo/bdjuno/modules/staking"
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
	cdc := ctx.EncodingConfig.Marshaler
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
	stakingModule := staking.NewModule(sources.StakingSource, slashingModule, cdc, db)
	rarimocoreModule := rarimocore.NewModule(sources.RarimoCoreSource, cdc, db)
	tokenmanagerModule := tokenmanager.NewModule(sources.TokenManagerSource, cdc, db)
	govModule := gov.NewModule(sources.GovSource, authModule, distrModule, mintModule, slashingModule, stakingModule, rarimocoreModule, cdc, db)

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
	}
}

package types

import (
	"fmt"
	rarimocoresource "gitlab.com/rarimo/bdjuno/modules/rarimocore/source"
	tokenmanagersource "gitlab.com/rarimo/bdjuno/modules/tokenmanager/source"
	rarimocoretypes "gitlab.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v3/node/remote"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/local"

	nodeconfig "github.com/forbole/juno/v3/node/config"

	banksource "gitlab.com/rarimo/bdjuno/modules/bank/source"
	localbanksource "gitlab.com/rarimo/bdjuno/modules/bank/source/local"
	remotebanksource "gitlab.com/rarimo/bdjuno/modules/bank/source/remote"
	distrsource "gitlab.com/rarimo/bdjuno/modules/distribution/source"
	localdistrsource "gitlab.com/rarimo/bdjuno/modules/distribution/source/local"
	remotedistrsource "gitlab.com/rarimo/bdjuno/modules/distribution/source/remote"
	govsource "gitlab.com/rarimo/bdjuno/modules/gov/source"
	localgovsource "gitlab.com/rarimo/bdjuno/modules/gov/source/local"
	remotegovsource "gitlab.com/rarimo/bdjuno/modules/gov/source/remote"
	mintsource "gitlab.com/rarimo/bdjuno/modules/mint/source"
	localmintsource "gitlab.com/rarimo/bdjuno/modules/mint/source/local"
	remotemintsource "gitlab.com/rarimo/bdjuno/modules/mint/source/remote"
	remoterarimocoresource "gitlab.com/rarimo/bdjuno/modules/rarimocore/source/remote"
	slashingsource "gitlab.com/rarimo/bdjuno/modules/slashing/source"
	localslashingsource "gitlab.com/rarimo/bdjuno/modules/slashing/source/local"
	remoteslashingsource "gitlab.com/rarimo/bdjuno/modules/slashing/source/remote"
	stakingsource "gitlab.com/rarimo/bdjuno/modules/staking/source"
	localstakingsource "gitlab.com/rarimo/bdjuno/modules/staking/source/local"
	remotestakingsource "gitlab.com/rarimo/bdjuno/modules/staking/source/remote"
	remotetokenmanagersource "gitlab.com/rarimo/bdjuno/modules/tokenmanager/source/remote"
)

type Sources struct {
	BankSource         banksource.Source
	DistrSource        distrsource.Source
	GovSource          govsource.Source
	MintSource         mintsource.Source
	SlashingSource     slashingsource.Source
	StakingSource      stakingsource.Source
	RarimoCoreSource   rarimocoresource.Source
	TokenManagerSource tokenmanagersource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	app := simapp.NewSimApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{},
	)

	sources := &Sources{
		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
		GovSource:      localgovsource.NewSource(source, govtypes.QueryServer(app.GovKeeper)),
		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(app, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(app, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(app, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:         remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:        remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:          remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		MintSource:         remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:     remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:      remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		RarimoCoreSource:   remoterarimocoresource.NewSource(source, rarimocoretypes.NewQueryClient(source.GrpcConn)),
		TokenManagerSource: remotetokenmanagersource.NewSource(source, tokenmanagertypes.NewQueryClient(source.GrpcConn)),
	}, nil
}

package types

import (
	"fmt"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	bridgesource "github.com/rarimo/bdjuno/modules/bridge/source"
	multisigsource "github.com/rarimo/bdjuno/modules/multisig/source"
	oraclemanagersource "github.com/rarimo/bdjuno/modules/oraclemanager/source"
	rarimocoresource "github.com/rarimo/bdjuno/modules/rarimocore/source"
	tokenmanagersource "github.com/rarimo/bdjuno/modules/tokenmanager/source"
	bridgetypes "github.com/rarimo/rarimo-core/x/bridge/types"
	multisigtypes "github.com/rarimo/rarimo-core/x/multisig/types"
	oraclemanagertypes "github.com/rarimo/rarimo-core/x/oraclemanager/types"
	rarimocoretypes "github.com/rarimo/rarimo-core/x/rarimocore/types"
	tokenmanagertypes "github.com/rarimo/rarimo-core/x/tokenmanager/types"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/node/remote"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v4/node/local"

	nodeconfig "github.com/forbole/juno/v4/node/config"

	banksource "github.com/rarimo/bdjuno/modules/bank/source"
	localbanksource "github.com/rarimo/bdjuno/modules/bank/source/local"
	remotebanksource "github.com/rarimo/bdjuno/modules/bank/source/remote"
	remotebridgesource "github.com/rarimo/bdjuno/modules/bridge/source/remote"
	distrsource "github.com/rarimo/bdjuno/modules/distribution/source"
	localdistrsource "github.com/rarimo/bdjuno/modules/distribution/source/local"
	remotedistrsource "github.com/rarimo/bdjuno/modules/distribution/source/remote"
	govsource "github.com/rarimo/bdjuno/modules/gov/source"
	localgovsource "github.com/rarimo/bdjuno/modules/gov/source/local"
	remotegovsource "github.com/rarimo/bdjuno/modules/gov/source/remote"
	mintsource "github.com/rarimo/bdjuno/modules/mint/source"
	localmintsource "github.com/rarimo/bdjuno/modules/mint/source/local"
	remotemintsource "github.com/rarimo/bdjuno/modules/mint/source/remote"
	remotemultisigsource "github.com/rarimo/bdjuno/modules/multisig/source/remote"
	remoteoraclemanagersource "github.com/rarimo/bdjuno/modules/oraclemanager/source/remote"
	remoterarimocoresource "github.com/rarimo/bdjuno/modules/rarimocore/source/remote"
	slashingsource "github.com/rarimo/bdjuno/modules/slashing/source"
	localslashingsource "github.com/rarimo/bdjuno/modules/slashing/source/local"
	remoteslashingsource "github.com/rarimo/bdjuno/modules/slashing/source/remote"
	stakingsource "github.com/rarimo/bdjuno/modules/staking/source"
	localstakingsource "github.com/rarimo/bdjuno/modules/staking/source/local"
	remotestakingsource "github.com/rarimo/bdjuno/modules/staking/source/remote"
	remotetokenmanagersource "github.com/rarimo/bdjuno/modules/tokenmanager/source/remote"
)

type Sources struct {
	BankSource          banksource.Source
	DistrSource         distrsource.Source
	GovSource           govsource.Source
	MintSource          mintsource.Source
	SlashingSource      slashingsource.Source
	StakingSource       stakingsource.Source
	RarimoCoreSource    rarimocoresource.Source
	TokenManagerSource  tokenmanagersource.Source
	OracleManagerSource oraclemanagersource.Source
	BridgeSource        bridgesource.Source
	MultisigSource      multisigsource.Source
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
		GovSource:      localgovsource.NewSource(source, govtypesv1.QueryServer(app.GovKeeper), nil),
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
		BankSource:          remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:         remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:           remotegovsource.NewSource(source, govtypesv1.NewQueryClient(source.GrpcConn), govtypesv1beta1.NewQueryClient(source.GrpcConn)),
		MintSource:          remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:      remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:       remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
		RarimoCoreSource:    remoterarimocoresource.NewSource(source, rarimocoretypes.NewQueryClient(source.GrpcConn)),
		TokenManagerSource:  remotetokenmanagersource.NewSource(source, tokenmanagertypes.NewQueryClient(source.GrpcConn)),
		OracleManagerSource: remoteoraclemanagersource.NewSource(source, oraclemanagertypes.NewQueryClient(source.GrpcConn)),
		BridgeSource:        remotebridgesource.NewSource(source, bridgetypes.NewQueryClient(source.GrpcConn)),
		MultisigSource:      remotemultisigsource.NewSource(source, multisigtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}

package gov

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"

	"gitlab.com/rarimo/bdjuno/types"
)

type ParamsUpdater interface {
	UpdateParams(height int64) error
}

type AuthModule interface {
	RefreshAccounts(height int64, addresses []string) error
}

type DistrModule = ParamsUpdater

type MintModule interface {
	ParamsUpdater
}

type SlashingModule = ParamsUpdater

type RarimoCoreModule interface {
	ParamsUpdater
	SaveOperationByIndex(height int64, index string) error
	HandleUpdateContract(height int64, details tokenmanagertypes.ContractUpgradeDetails) error
	GetFeeToken(height int64, chain, contract string) (*tokenmanagertypes.FeeToken, error)
}

type BridgeModule = ParamsUpdater

type OracleManagerModule interface {
	ParamsUpdater
	HandleOracle(height int64, chain, account string) error
}

type TokenManagerModule interface {
	ParamsUpdater
	UpdateItems(items []*tokenmanagertypes.Item) error
	RemoveItems(indexes []string) error
	CreateCollection(
		index string,
		meta *tokenmanagertypes.CollectionMetadata,
		data []*tokenmanagertypes.CollectionData,
		items []*tokenmanagertypes.Item,
		onChainItems []*tokenmanagertypes.OnChainItem,
	) error
	UpdateCollectionDatas(datas []*tokenmanagertypes.CollectionData) error
	CreateCollectionDatas(height int64, datas []*tokenmanagertypes.CollectionData) error
	RemoveCollectionDatas(height int64, indexes []*tokenmanagertypes.CollectionDataIndex) error
	RemoveCollection(index string) error
}

type StakingModule interface {
	ParamsUpdater
	GetStakingPool(height int64) (*types.Pool, error)
	GetStakingPoolSnapshot(height int64) (*types.PoolSnapshot, error)
	GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, []types.Validator, error)
	GetValidatorsVotingPowers(height int64, vals *tmctypes.ResultValidators) ([]types.ValidatorVotingPower, error)
	GetValidatorsStatuses(height int64, validators []stakingtypes.Validator) ([]types.ValidatorStatus, error)
}

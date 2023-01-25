package gov

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tokenmanagertypes "gitlab.com/rarimo/rarimo-core/x/tokenmanager/types"

	"gitlab.com/rarimo/bdjuno/types"
)

type AuthModule interface {
	RefreshAccounts(height int64, addresses []string) error
}

type DistrModule interface {
	UpdateParams(height int64) error
}

type MintModule interface {
	UpdateParams(height int64) error
	UpdateInflation() error
}

type SlashingModule interface {
	UpdateParams(height int64) error
}

type RarimoCoreModule interface {
	UpdateParams(height int64) error
}

type TokenManagerModule interface {
	UpdateParams(height int64) error
	UpdateItems(items []*tokenmanagertypes.Item) error
	RemoveItems(indexes []*tokenmanagertypes.ItemIndex) error
	CreateCollection(index string, meta *tokenmanagertypes.CollectionMetadata, data []*tokenmanagertypes.CollectionData) error
	UpdateCollectionDatas(datas []*tokenmanagertypes.CollectionData) error
	CreateCollectionDatas(datas []*tokenmanagertypes.CollectionData) error
	RemoveCollectionDatas(indexes []*tokenmanagertypes.CollectionDataIndex) error
	RemoveCollection(index string) error
}

type StakingModule interface {
	GetStakingPool(height int64) (*types.Pool, error)
	GetValidatorsWithStatus(height int64, status string) ([]stakingtypes.Validator, []types.Validator, error)
	GetValidatorsVotingPowers(height int64, vals *tmctypes.ResultValidators) ([]types.ValidatorVotingPower, error)
	GetValidatorsStatuses(height int64, validators []stakingtypes.Validator) ([]types.ValidatorStatus, error)
	UpdateParams(height int64) error
}

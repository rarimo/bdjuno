package local

import (
	"fmt"
	juno "github.com/forbole/juno/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v4/node/local"

	"github.com/rarimo/bdjuno/modules/bank/source"
	"github.com/rarimo/bdjuno/types"
)

var (
	_      source.Source = &Source{}
	events               = map[string]interface{}{
		"coin_spent":    nil,
		"coin_received": nil,
	}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q banktypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk banktypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      bk,
	}
}

// GetBalances implements keeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var balances []types.AccountBalance
	for _, address := range addresses {
		res, err := s.q.AllBalances(sdk.WrapSDKContext(ctx), &banktypes.QueryAllBalancesRequest{Address: address})
		if err != nil {
			return nil, err
		}

		balances = append(balances, types.NewAccountBalance(address, res.Balances, height))
	}

	return balances, nil
}

// GetSupply implements bankkeeper.Source
func (s Source) GetSupply(height int64) (sdk.Coins, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var coins []sdk.Coin
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.q.TotalSupply(
			sdk.WrapSDKContext(ctx),
			&banktypes.QueryTotalSupplyRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 supplies at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting total supply: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		coins = append(coins, res.Supply...)
	}

	return coins, nil
}

// GetAccountBalances implements bankkeeper.Source
func (s Source) GetAccountBalance(address string, height int64) ([]sdk.Coin, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	balRes, err := s.q.AllBalances(sdk.WrapSDKContext(ctx), &banktypes.QueryAllBalancesRequest{Address: address})
	if err != nil {
		return nil, fmt.Errorf("error while getting all balances: %s", err)
	}

	return balRes.Balances, nil
}

// GetEventsFromTxs implements bankkeeper.Source
func (s Source) GetEventsFromTxs(txs []*juno.Tx) ([]source.TxEvent, error) {
	res := make([]source.TxEvent, 0)

	var commonError error
	for _, tx := range txs {
		for _, event := range tx.Events {
			if _, ok := events[event.Type]; ok {
				coin, err := sdk.ParseCoinNormalized(string(event.Attributes[1].Value))
				if err != nil {
					log.Err(err).Msg("failed to parse coins")
					commonError = errors.Wrap(err, "failed to parse coins")
					continue
				}

				res = append(res, source.TxEvent{
					BlockHeight: tx.Height,
					Address:     string(event.Attributes[0].Value),
					Event:       event.Type,
					Amount:      sdk.NewCoins(coin),
					TxHash:      tx.TxHash,
				})
			}
		}
	}

	return res, commonError
}

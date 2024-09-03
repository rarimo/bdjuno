package remote

import (
	"fmt"
	juno "github.com/forbole/juno/v4/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/juno/v4/node/remote"

	bankkeeper "github.com/rarimo/bdjuno/modules/bank/source"
	"github.com/rarimo/bdjuno/types"
)

var (
	_      bankkeeper.Source = &Source{}
	events                   = map[string]interface{}{
		"coin_spent":    nil,
		"coin_received": nil,
	}
)

type Source struct {
	*remote.Source
	bankClient banktypes.QueryClient
	client     *tmclient.HTTP
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, bankClient banktypes.QueryClient) *Source {
	return &Source{
		Source:     source,
		bankClient: bankClient,
	}
}

// GetBalances implements bankkeeper.Source
func (s Source) GetBalances(addresses []string, height int64) ([]types.AccountBalance, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var balances []types.AccountBalance
	for _, address := range addresses {
		balRes, err := s.bankClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: address})
		if err != nil {
			return nil, fmt.Errorf("error while getting all balances: %s", err)
		}

		balances = append(balances, types.NewAccountBalance(
			address,
			balRes.Balances,
			height,
		))
	}

	return balances, nil
}

// GetSupply implements bankkeeper.Source
func (s Source) GetSupply(height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var coins []sdk.Coin
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.bankClient.TotalSupply(
			ctx,
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

// GetEventsFromTxs implements bankkeeper.Source
func (s Source) GetEventsFromTxs(txs []*juno.Tx) ([]bankkeeper.TxEvent, error) {
	res := make([]bankkeeper.TxEvent, 0)

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

				res = append(res, bankkeeper.TxEvent{
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

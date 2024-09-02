package remote

import (
	"context"
	"fmt"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	"log"

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
func NewSource(source *remote.Source, bankClient banktypes.QueryClient, rpcURL string) *Source {
	client, err := tmclient.New(rpcURL, "")
	if err != nil {
		log.Fatalf("Failed to create Tendermint RPC client: %v", err)
	}

	return &Source{
		Source:     source,
		bankClient: bankClient,
		client:     client,
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

func (s Source) GetEventsByBlockHeight(height int64) []bankkeeper.TxEvent {
	// Get block results which include events
	blockResults, err := s.client.BlockResults(context.Background(), &height)
	if err != nil {
		log.Fatalf("Failed to get block results: %v", err)
	}

	res := make([]bankkeeper.TxEvent, 0)

	if blockResults.TxsResults == nil {
		return nil
	}
	// Iterate over the transactions in the block and print the events
	for _, txResult := range blockResults.TxsResults {
		for _, event := range txResult.Events {
			if _, ok := events[event.Type]; ok {
				coin, err := sdk.ParseCoinNormalized(string(event.Attributes[1].Value))
				if err != nil {
					log.Println("failed to parse coins: ", err)
					continue
				}

				res = append(res, bankkeeper.TxEvent{
					BlockHeight: height,
					Address:     string(event.Attributes[0].Value),
					Event:       event.Type,
					Amount:      sdk.NewCoins(coin),
				})
			}

		}
	}
	return res
}

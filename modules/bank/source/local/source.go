package local

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/node/local"

	"github.com/rarimo/bdjuno/modules/bank/source"
	"github.com/rarimo/bdjuno/types"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the bank keeper that works on a local node
type Source struct {
	*local.Source
	q      banktypes.QueryServer
	client *tmclient.HTTP
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk banktypes.QueryServer, rpcURL string) *Source {
	client, err := tmclient.New(rpcURL, "")
	if err != nil {
		log.Fatalf("Failed to create Tendermint RPC client: %v", err)
	}

	return &Source{
		Source: source,
		q:      bk,
		client: client,
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

func (s Source) GetEventsByBlockHeight(height int64) []source.TxEvent {
	// Get block results which include events
	blockResults, err := s.client.BlockResults(context.Background(), &height)
	if err != nil {
		log.Fatalf("Failed to get block results: %v", err)
	}

	res := make([]source.TxEvent, 0)

	// Iterate over the transactions in the block and print the events
	for _, txResult := range blockResults.TxsResults {
		for _, event := range txResult.Events {
			fmt.Println("event.String()", event.String())

			for _, attribute := range event.Attributes {
				s.Logger.Debug(" - %s: %s\n", string(attribute.Key), string(attribute.Value))
			}
		}
	}
	return res
}

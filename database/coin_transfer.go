package database

import (
	"fmt"

	dbtypes "github.com/rarimo/bdjuno/database/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
)

// SaveSupply allows to save for the given height the given total amount of coins
func (db *Db) SaveCoinTransfer(coin sdk.Coins, height int64, event, address string) error {
	query := `
		INSERT INTO coin_transfer (coins, height, event, address) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.SQL.Exec(query, pq.Array(dbtypes.NewDbCoins(coin)), height, event, address)
	if err != nil {
		return fmt.Errorf("error while storing supply: %s", err)
	}

	return nil
}

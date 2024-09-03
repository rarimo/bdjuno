package database

import (
	"fmt"

	dbtypes "github.com/rarimo/bdjuno/database/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
)

// SaveCoinTransfer allows to save tokens transfers
func (db *Db) SaveCoinTransfer(coin sdk.Coins, height int64, event, address string, txHash string) error {
	query := `
		INSERT INTO coin_transfer (coins, height, event, address, tx_hash) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.SQL.Exec(query, pq.Array(dbtypes.NewDbCoins(coin)), height, event, address, txHash)
	if err != nil {
		return fmt.Errorf("error while storing tokens transfer: %s", err)
	}

	return nil
}

package types

// CoinTransferRow represents a coin transfer
type CoinTransferRow struct {
	ID      int64   `db:"id"`
	Coin    *DbCoin `db:"coin"`
	Height  int64   `db:"height"`
	Event   string  `db:"event"`
	Address string  `db:"address"`
	TxHash  string  `db:"tx_hash"`
}

// NewCoinTransferRow allows to easily create a new CoinTransferRow
func NewCoinTransferRow(id int64, coin DbCoin, height int64, event, address string, txHash string) CoinTransferRow {
	return CoinTransferRow{
		ID:      id,
		Coin:    &coin,
		Height:  height,
		Event:   event,
		Address: address,
		TxHash:  txHash,
	}
}

package types

// SupplyRow represents a single row inside the "supply" table
type CoinTransferRow struct {
	ID      int64   `db:"id"`
	Coin    *DbCoin `db:"coin"`
	Height  int64   `db:"height"`
	Event   string  `db:"event"`
	Address string  `db:"address"`
}

// NewSupplyRow allows to easily create a new NewSupplyRow
func NewCoinTransferRow(id int64, coin DbCoin, height int64, event, address string) CoinTransferRow {
	return CoinTransferRow{
		ID:      id,
		Coin:    &coin,
		Height:  height,
		Event:   event,
		Address: address,
	}
}

package types

// ItemRow represents a single row of the "item" table
type ItemRow struct {
	Index      string `db:"index"`
	Collection string `db:"collection"`
	Meta       string `db:"meta"`
	OnChain    string `db:"on_chain"`
}

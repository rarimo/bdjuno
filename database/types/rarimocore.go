package types

// OperationRow represents a single row of the "operation" table
type OperationRow struct {
	Index         string `db:"index"`
	OperationType int32  `db:"operation_type"`
	Status        int32  `db:"status"`
	Creator       string `db:"creator"`
	Timestamp     uint64 `db:"timestamp"`
}

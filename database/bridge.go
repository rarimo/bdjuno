package database

import (
	"fmt"
	"gitlab.com/rarimo/bdjuno/types"
	"strings"
)

// SaveBridgeParams saves the given x/bridge parameters inside the database
func (db *Db) SaveBridgeParams(params *types.BridgeParams) (err error) {
	stmt := `
INSERT INTO bridge_params(withdraw_denom, height)
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE
	SET withdraw_denom = excluded.withdraw_denom,
		height = excluded.height
WHERE bridge_params.height <= excluded.height
`
	_, err = db.Sql.Exec(
		stmt,
		params.WithdrawDenom,
		params.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing bridge params: %s", err)
	}

	return nil
}

func (db *Db) SaveHashes(hashes []types.Hash) error {
	if len(hashes) == 0 {
		return nil
	}

	query := `INSERT INTO hash (index) VALUES `

	var params []interface{}

	for i, hash := range hashes {
		// Prepare the hash query
		vi := i * 1
		query += fmt.Sprintf("($%d),", vi+1)
		params = append(params, hash.Index)
	}

	// Store the hashes
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing hashes: %s", err)
	}

	return nil
}

package database

import (
	"fmt"
	"gitlab.com/rarimo/bdjuno/types"
	"strings"
)

// SaveOracleManagerParams saves the given x/oraclemanager parameters inside the database
func (db *Db) SaveOracleManagerParams(params *types.OracleManagerParams) (err error) {
	stmt := `
INSERT INTO oraclemanager_params(min_oracle_stake, check_operation_delta, max_violations_count, max_missed_count, slashed_freeze_blocks, min_oracles_count, stake_denom, vote_quorum, vote_threshold, height)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (one_row_id) DO UPDATE
	SET min_oracle_stake = excluded.min_oracle_stake,
		check_operation_delta = excluded.check_operation_delta,
		max_violations_count = excluded.max_violations_count,
		max_missed_count = excluded.max_missed_count,
		slashed_freeze_blocks = excluded.slashed_freeze_blocks,
		min_oracles_count = excluded.min_oracles_count,
		stake_denom = excluded.stake_denom,
		vote_quorum = excluded.vote_quorum,
		vote_threshold = excluded.vote_threshold,
		height = excluded.height
WHERE oraclemanager_params.height <= excluded.height
`
	_, err = db.Sql.Exec(
		stmt,
		params.MinOracleStake,
		params.CheckOperationDelta,
		params.MaxViolationsCount,
		params.MaxMissedCount,
		params.SlashedFreezeBlocks,
		params.MinOraclesCount,
		params.StakeDenom,
		params.VoteQuorum,
		params.VoteThreshold,
		params.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing oraclemanager params: %s", err)
	}

	return nil
}

func (db *Db) SaveOracles(oracles []types.Oracle) error {
	if len(oracles) == 0 {
		return nil
	}

	var accounts []types.Account

	query := `INSERT INTO oracle (index, chain, account, status, stake, missed_count, violations_count, freeze_end_block, votes_count, create_operations_count) VALUES `

	var params []interface{}

	for i, oracle := range oracles {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(oracle.Account))

		// Prepare the oracle query
		vi := i * 10
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10)

		params = append(
			params,
			oracle.Index,
			oracle.Chain,
			oracle.Account,
			oracle.Status,
			oracle.Stake,
			oracle.MissedCount,
			oracle.ViolationsCount,
			oracle.FreezeEndBlock,
			oracle.VotesCount,
			oracle.CreateOperationsCount,
		)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing oracle accounts: %s", err)
	}

	// Store the oracles
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += ` ON CONFLICT (index) DO UPDATE 
	SET status = excluded.status, stake = excluded.stake, violations_count = excluded.violations_count, freeze_end_block = excluded.freeze_end_block, votes_count = excluded.votes_count, create_operations_count = excluded.create_operations_count
WHERE oracle.index = excluded.index
`
	_, err = db.Sql.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing oracles: %s", err)
	}

	return nil
}

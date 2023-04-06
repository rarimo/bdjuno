package database

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	dbtypes "gitlab.com/rarimo/bdjuno/database/types"
	"gitlab.com/rarimo/bdjuno/types"
	"strings"
)

// SaveParties saves the given x/gov parameters inside the database
func (db *Db) SaveParties(parties []types.Party) error {
	if len(parties) == 0 {
		return nil
	}

	var accounts []types.Account

	partiesQuery := `INSERT INTO parties (address, account, pub_key, status, violations_count, freeze_end_block, delegator) VALUES `

	var partiesParams []interface{}

	for i, party := range parties {
		accounts = append(accounts, types.NewAccount(party.Account))

		vi := i * 7
		partiesQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7)

		partiesParams = append(partiesParams, party.Address, party.Account, party.PubKey, party.Status, party.ViolationsCount, party.FreezeEndBlock, party.Delegator)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing proposers accounts: %s", err)
	}

	// Store the proposals
	partiesQuery = strings.TrimSuffix(partiesQuery, ",") // Remove trailing ","
	partiesQuery += ` ON CONFLICT (account) DO UPDATE 
	SET status = excluded.status, pub_key = excluded.pub_key, violations_count = excluded.violations_count, freeze_end_block = excluded.freeze_end_block, delegator = excluded.delegator
WHERE parties.account = excluded.account
`
	_, err = db.Sql.Exec(partiesQuery, partiesParams...)
	if err != nil {
		return fmt.Errorf("error while storing parties: %s", err)
	}

	return nil
}

// SaveRarimoCoreParams saves the given x/rarimocore parameters inside the database
func (db *Db) SaveRarimoCoreParams(params *types.RarimoCoreParams) (err error) {
	stmt := `
INSERT INTO rarimocore_params(key_ecdsa, threshold, is_update_required, last_signature, stake_amount, stake_denom, max_violations_count, freeze_blocks_period, parties, height)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (one_row_id) DO UPDATE
	SET key_ecdsa = excluded.key_ecdsa,
		threshold = excluded.threshold,
		is_update_required = excluded.is_update_required,
		last_signature = excluded.last_signature,
		stake_amount = excluded.stake_amount,
		stake_denom = excluded.stake_denom,
		max_violations_count = excluded.max_violations_count,
		freeze_blocks_period = excluded.freeze_blocks_period,
		parties = excluded.parties,
		height = excluded.height
WHERE rarimocore_params.height <= excluded.height
`
	_, err = db.Sql.Exec(
		stmt,
		params.KeyECDSA,
		params.Threshold,
		params.IsUpdateRequired,
		params.LastSignature,
		params.StakeAmount,
		params.StakeDenom,
		params.MaxViolationsCount,
		params.FreezeBlocksPeriod,
		pq.StringArray(params.Parties),
		params.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing rarimocore params: %s", err)
	}

	return nil
}

func (db *Db) SaveOperations(operations []types.Operation) error {
	if len(operations) == 0 {
		return nil
	}

	var accounts []types.Account

	operationsQuery := `INSERT INTO operation (index, operation_type, status, creator, timestamp) VALUES `

	var operationsParams []interface{}

	for i, operation := range operations {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(operation.Creator))

		// Prepare the operation query
		vi := i * 5
		operationsQuery += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5)

		operationsParams = append(
			operationsParams,
			operation.Index,
			operation.OperationType,
			operation.Status,
			operation.Creator,
			operation.Timestamp,
		)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing operation accounts: %s", err)
	}

	// Store the operations
	operationsQuery = strings.TrimSuffix(operationsQuery, ",") // Remove trailing ","
	operationsQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(operationsQuery, operationsParams...)
	if err != nil {
		return fmt.Errorf("error while storing operations: %s", err)
	}

	return nil
}

func (db *Db) UpdateOperation(operation types.Operation) error {
	query := `UPDATE operation SET status = $1 WHERE index = $2`
	_, err := db.Sql.Exec(query,
		operation.Status,
		operation.Index,
	)
	if err != nil {
		return fmt.Errorf("error while updating operation: %s", err)
	}

	return nil
}

func (db *Db) GetOperation(index string) (*types.Operation, error) {
	var rows []*dbtypes.OperationRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM operation WHERE index = $1`, index)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	row := rows[0]

	operation := types.NewOperation(
		row.Index,
		row.OperationType,
		row.Status,
		row.Creator,
		row.Timestamp,
	)

	return &operation, nil
}

func (db *Db) SaveTransfers(transfers []types.Transfer) (err error) {
	if len(transfers) == 0 {
		return nil
	}

	transfersQuery := `
INSERT INTO transfer (
	operation_index, origin, tx, event_id, "from", "to", receiver, amount, bundle_data, 
    bundle_salt, item_meta
) VALUES`

	var transfersParams []interface{}

	for i, transfer := range transfers {
		// Prepare the transfer query
		vi := i * 11
		transfersQuery += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11,
		)

		from, err := json.Marshal(transfer.From)
		if err != nil {
			return fmt.Errorf("error while unmarshaling transfer.From: %s", err)
		}

		to, err := json.Marshal(transfer.To)
		if err != nil {
			return fmt.Errorf("error while unmarshaling transfer.To: %s", err)
		}

		meta, err := json.Marshal(transfer.ItemMeta)
		if err != nil {
			return fmt.Errorf("error while unmarshaling transfer.ItemMeta: %s", err)
		}

		transfersParams = append(transfersParams,
			transfer.OperationIndex,
			transfer.Origin,
			transfer.Tx,
			transfer.EventID,
			string(from),
			string(to),
			transfer.Receiver,
			transfer.Amount,
			transfer.BundleData,
			transfer.BundleSalt,
			string(meta),
		)
	}

	transfersQuery = strings.TrimSuffix(transfersQuery, ",") // Remove trailing ","
	transfersQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(transfersQuery, transfersParams...)
	if err != nil {
		return fmt.Errorf("error while storing transfers: %s", err)
	}

	return nil
}

func (db *Db) SaveChangeParties(changeParties []types.ChangeParties) (err error) {
	if len(changeParties) == 0 {
		return nil
	}

	changePartiesQuery := `INSERT INTO change_parties (operation_index, parties, new_public_key, signature) VALUES`
	var changePartiesParams []interface{}

	for i, changeParty := range changeParties {
		vi := i * 4
		changePartiesQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4)

		changePartiesParams = append(changePartiesParams,
			changeParty.OperationIndex,
			pq.StringArray(changeParty.Parties),
			changeParty.NewPublicKey,
			changeParty.Signature,
		)
	}

	changePartiesQuery = strings.TrimSuffix(changePartiesQuery, ",") // Remove trailing ","
	changePartiesQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(changePartiesQuery, changePartiesParams...)
	if err != nil {
		return fmt.Errorf("error while storing change parties: %s", err)
	}

	return nil
}

func (db *Db) UpdateChangeParties(changeParties types.ChangeParties) (err error) {
	query := `UPDATE change_parties SET parties = $1, new_public_key = $2, signature = $3 WHERE operation_index = $4`
	_, err = db.Sql.Exec(query,
		pq.StringArray(changeParties.Parties),
		changeParties.NewPublicKey,
		changeParties.Signature,
		changeParties.OperationIndex,
	)
	if err != nil {
		return fmt.Errorf("error while updating change parties: %s", err)
	}

	return nil
}

func (db *Db) SaveConfirmations(confirmations []types.Confirmation) (err error) {
	if len(confirmations) == 0 {
		return nil
	}

	confirmationsQuery := `INSERT INTO confirmation (root, indexes, signature_ecdsa, creator) VALUES`
	var confirmationsParams []interface{}

	for i, confirmation := range confirmations {
		vi := i * 4
		confirmationsQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4)

		confirmationsParams = append(confirmationsParams,
			confirmation.Root,
			pq.StringArray(confirmation.Indexes),
			confirmation.SignatureECDSA,
			confirmation.Creator,
		)
	}

	confirmationsQuery = strings.TrimSuffix(confirmationsQuery, ",") // Remove trailing ","
	confirmationsQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(confirmationsQuery, confirmationsParams...)
	if err != nil {
		return fmt.Errorf("error while storing confirmations: %s", err)
	}

	return nil
}

func (db *Db) SaveRarimoCoreVotes(votes []types.RarimoCoreVote) (err error) {
	if len(votes) == 0 {
		return nil
	}

	query := `INSERT INTO vote (operation, validator, vote) VALUES `
	var queryParams []interface{}

	for i, vote := range votes {
		vi := i * 3
		query += fmt.Sprintf("($%d, $%d, $%d),", vi+1, vi+2, vi+3)
		queryParams = append(queryParams, vote.Operation, vote.Validator, vote.Vote)
	}

	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"

	_, err = db.Sql.Exec(query, queryParams...)
	if err != nil {
		return fmt.Errorf("error while storing votes: %s", err)
	}

	return nil
}

func (db *Db) RemoveRarimoCoreVotes(opIndex string) error {
	stmt := `DELETE FROM vote WHERE operation = $1`
	_, err := db.Sql.Exec(stmt, opIndex)
	if err != nil {
		return fmt.Errorf("error while deleting rarimo core votes: %s", err)
	}

	return nil
}

func (db *Db) SaveViolationReports(reports []types.ViolationReport) (err error) {
	if len(reports) == 0 {
		return nil
	}

	query := `INSERT INTO violation_report (index, session_id, offender, sender, violation_type, msg) VALUES `
	var queryParams []interface{}

	for i, report := range reports {
		vi := i * 6
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6)
		queryParams = append(queryParams, report.Index, report.SessionId, report.Offender, report.Sender, report.ViolationType, report.Msg)
	}

	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"

	_, err = db.Sql.Exec(query, queryParams...)
	if err != nil {
		return fmt.Errorf("error while storing reports: %s", err)
	}

	return nil
}

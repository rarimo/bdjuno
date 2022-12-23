package database

import (
	"fmt"
	dbtypes "gitlab.com/rarimo/bdjuno/database/types"
	"gitlab.com/rarimo/bdjuno/types"
)

// SaveParties saves the given x/gov parameters inside the database
func (db *Db) SaveParties(parties []types.Party) error {
	if len(parties) == 0 {
		return nil
	}

	var accounts []types.Account

	partiesQuery := `INSERT INTO party (address, account, public_key, verified) VALUES `

	var partiesParams []interface{}

	for i, party := range parties {
		accounts = append(accounts, types.NewAccount(party.Account))

		vi := i * 4
		partiesQuery += fmt.Sprintf("($%d, $%d, $%d, $%d),", vi+1, vi+2, vi+3, vi+4)

		partiesParams = append(partiesParams, party.Address, party.Account, party.PubKey, party.Verified)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing proposers accounts: %s", err)
	}

	// Store the proposals
	partiesQuery = partiesQuery[:len(partiesQuery)-1] // Remove trailing ","
	partiesQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(partiesQuery, partiesParams...)
	if err != nil {
		return fmt.Errorf("error while storing parties: %s", err)
	}

	return nil
}

// SaveRarimoCoreParams saves the given x/rarimocore parameters inside the database
func (db *Db) SaveRarimoCoreParams(params *types.RarimoCoreParams) (err error) {
	stmt := `
INSERT INTO rarimocore_params(key_ecdsa, threshold, is_update_required, last_signature, parties, height)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (one_row_id) DO UPDATE
	SET key_ecdsa = excluded.key_ecdsa,
		threshold = excluded.threshold,
		is_update_required = excluded.is_update_required,
		last_signature = excluded.last_signature,
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
		params.Parties,
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

	operationsQuery := `INSERT INTO operation (index, operation_type, signed, creator, timestamp) VALUES `

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
			operation.Signed,
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
	operationsQuery = operationsQuery[:len(operationsQuery)-1] // Remove trailing ","
	operationsQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(operationsQuery, operationsParams...)
	if err != nil {
		return fmt.Errorf("error while storing operations: %s", err)
	}

	return nil
}

func (db *Db) UpdateOperation(operation types.Operation) error {
	query := `UPDATE operation SET signed = $1, where index = $2`
	_, err := db.Sql.Exec(query,
		operation.Signed,
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
		row.Signed,
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
	operation_index, origin, tx, event_id, from_chain, to_chain, receiver, amount, bundle_data, 
    bundle_salt, token_index
) VALUES`

	var transfersParams []interface{}

	for i, transfer := range transfers {
		// Prepare the transfer query
		vi := i * 11
		transfersQuery += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11,
		)

		transfersParams = append(transfersParams,
			transfer.OperationIndex,
			transfer.Origin,
			transfer.Tx,
			transfer.EventID,
			transfer.FromChain,
			transfer.ToChain,
			transfer.Receiver,
			transfer.Amount,
			transfer.BundleData,
			transfer.BundleSalt,
			transfer.TokenIndex,
		)
	}

	transfersQuery = transfersQuery[:len(transfersQuery)-1] // Remove trailing ","
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
			changeParty.Parties,
			changeParty.NewPublicKey,
			changeParty.Signature,
		)
	}

	changePartiesQuery = changePartiesQuery[:len(changePartiesQuery)-1] // Remove trailing ","
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
		changeParties.Parties,
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
			confirmation.Indexes,
			confirmation.SignatureECDSA,
			confirmation.Creator,
		)
	}

	confirmationsQuery = confirmationsQuery[:len(confirmationsQuery)-1] // Remove trailing ","
	confirmationsQuery += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(confirmationsQuery, confirmationsParams...)
	if err != nil {
		return fmt.Errorf("error while storing confirmations: %s", err)
	}

	return nil
}

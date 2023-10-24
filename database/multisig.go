package database

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/rarimo/bdjuno/types"
	"strings"
)

// SaveMultisigParams saves the given x/multisig parameters inside the database
func (db *Db) SaveMultisigParams(params *types.MultisigParams) (err error) {
	stmt := `
INSERT INTO multisig_params(group_sequence, proposal_sequence, prune_period, voting_period, height)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (one_row_id) DO UPDATE
	SET group_sequence = excluded.group_sequence,
		proposal_sequence = excluded.proposal_sequence,
		prune_period = excluded.prune_period,
		voting_period = excluded.voting_period,
		height = excluded.height
WHERE multisig_params.height <= excluded.height
`
	_, err = db.SQL.Exec(
		stmt,
		params.GroupSequence,
		params.ProposalSequence,
		params.PrunePeriod,
		params.VotingPeriod,
		params.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing multisig params: %s", err)
	}

	return nil
}

// SaveGroups saves the given x/multisig groups inside the database
func (db *Db) SaveGroups(groups []*types.Group) (err error) {
	if len(groups) == 0 {
		return nil
	}

	var accounts []types.Account

	query := `
INSERT INTO "group"(
	account, members, threshold
) VALUES`

	var params []interface{}

	for i, group := range groups {
		membersAccounts := make([]types.Account, len(group.Members))
		for j, member := range group.Members {
			membersAccounts[j] = types.NewAccount(member)
		}

		accounts = append(accounts, types.NewAccount(group.Account))
		accounts = append(accounts, membersAccounts...)

		vi := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)
		params = append(params, group.Account, pq.StringArray(group.Members), group.Threshold)
	}

	// Store the groups
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += ` ON CONFLICT (account) DO UPDATE
 	SET members = excluded.members, threshold = excluded.threshold
WHERE "group".account = excluded.account
 `

	err = db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing group accounts: %s", err)
	}

	_, err = db.SQL.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing multisig groups: %s", err)
	}

	return nil
}

// SaveMultisigProposals saves the given x/multisig proposals inside the database
func (db *Db) SaveMultisigProposals(proposals []*types.MultisigProposal) (err error) {
	if len(proposals) == 0 {
		return nil
	}

	var accounts []types.Account

	query := `
INSERT INTO multisig_proposal(
	id, proposer, "group", submit_block, voting_end_block, status, final_tally_result, messages
) VALUES`

	var params []interface{}

	for i, proposal := range proposals {
		accounts = append(accounts, types.NewAccount(proposal.Proposer))

		vi := i * 8
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8)

		messages := "["
		for j, message := range proposal.Messages {
			msg, err := db.EncodingConfig.Codec.MarshalJSON(message)
			if err != nil {
				return fmt.Errorf("error while marshaling message: %s", err)
			}

			messages += string(msg)
			if j != len(proposal.Messages)-1 {
				messages += ","
			}

		}

		messages += "]"
		result := json.RawMessage("{}")

		if proposal.FinalTallyResult != nil {
			result, err = json.Marshal(proposal.FinalTallyResult)
			if err != nil {
				return fmt.Errorf("error while marshaling final tally result: %s", err)
			}
		}

		params = append(params, proposal.Id, proposal.Proposer, proposal.Group, proposal.SubmitBlock, proposal.VotingEndBlock, proposal.Status, string(result), messages)
	}

	err = db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing multisig proposers accounts: %s", err)
	}

	// Store the proposals
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += ` ON CONFLICT (id) DO UPDATE
 	SET final_tally_result = excluded.final_tally_result, status = excluded.status
WHERE multisig_proposal.id = excluded.id
 `
	_, err = db.SQL.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing multisig proposals: %s", err)
	}

	return nil
}

// SaveMultisigProposalVotes saves the given x/multisig proposal votes inside the database
func (db *Db) SaveMultisigProposalVotes(votes []*types.MultisigProposalVote) (err error) {
	if len(votes) == 0 {
		return nil
	}

	query := `
INSERT INTO multisig_proposal_vote(
	index, voter, proposal_id, option, submit_block
) VALUES`

	var params []interface{}

	for i, vote := range votes {
		vi := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5)
		params = append(params, vote.Index, vote.Voter, vote.ProposalId, vote.Option, vote.SubmitBlock)
	}

	// Store the votes
	query = strings.TrimSuffix(query, ",") // Remove trailing ","
	query += ` ON CONFLICT (index) DO UPDATE
 	SET option = excluded.option, submit_block = excluded.submit_block
WHERE multisig_proposal_vote.index = excluded.index
 `
	_, err = db.SQL.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("error while storing multisig votes: %s", err)
	}

	return nil
}

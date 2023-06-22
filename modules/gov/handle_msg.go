package gov

import (
	"fmt"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"strconv"

	"gitlab.com/rarimo/bdjuno/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	juno "github.com/forbole/juno/v4/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *govtypesv1beta1.MsgSubmitProposal:
		return m.handleMsgSubmitProposal(tx, index, cosmosMsg)

	case *govtypesv1.MsgDeposit:
		return m.handleMsgDeposit(tx, cosmosMsg)

	case *govtypesv1.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func (m *Module) handleMsgSubmitProposal(tx *juno.Tx, index int, msg *govtypesv1beta1.MsgSubmitProposal) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(tx.Height, proposalID)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	// Unpack the content
	var content govtypesv1beta1.Content
	err = m.cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while unpacking proposal content: %s", err)
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.ProposalId,
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.GetContent(),
		proposal.Status.String(),
		proposal.SubmitBlock,
		proposal.DepositEndBlock,
		proposal.VotingStartBlock,
		proposal.VotingEndBlock,
		msg.Proposer,
	)
	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func (m *Module) handleMsgDeposit(tx *juno.Tx, msg *govtypesv1.MsgDeposit) error {
	deposit, err := m.source.ProposalDeposit(tx.Height, msg.ProposalId, msg.Depositor)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	return m.db.SaveDeposits([]types.Deposit{
		types.NewDeposit(msg.ProposalId, msg.Depositor, deposit.Amount, tx.Height),
	})
}

// handleMsgVote allows to properly handle a handleMsgVote
func (m *Module) handleMsgVote(tx *juno.Tx, msg *govtypesv1.MsgVote) error {
	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, tx.Height)
	return m.db.SaveVote(vote)
}

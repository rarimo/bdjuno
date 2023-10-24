package gov

import (
	"encoding/json"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"
	"strconv"
	"strings"

	"github.com/rarimo/bdjuno/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
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

	case *govtypesv1.MsgSubmitProposal:
		return m.handleMsgSubmitProposalV1(tx, index, cosmosMsg)

	case *govtypesv1beta1.MsgDeposit:
		return m.handleMsgDeposit(tx.Height, cosmosMsg.ProposalId, cosmosMsg.Depositor)

	case *govtypesv1.MsgDeposit:
		return m.handleMsgDeposit(tx.Height, cosmosMsg.ProposalId, cosmosMsg.Depositor)

	case *govtypesv1beta1.MsgVote:
		return m.handleMsgVote(tx.Height, cosmosMsg.ProposalId, cosmosMsg.Voter, int32(cosmosMsg.Option))

	case *govtypesv1.MsgVote:
		return m.handleMsgVote(tx.Height, cosmosMsg.ProposalId, cosmosMsg.Voter, int32(cosmosMsg.Option))
	}

	return nil
}

func (m *Module) handleMsgSubmitProposalV1(tx *juno.Tx, index int, msg *govtypesv1.MsgSubmitProposal) error {
	proposalID, err := m.getProposalID(tx, index)
	if err != nil {
		return fmt.Errorf("error while getting proposal id: %s", err)
	}

	proposal, err := m.source.ProposalV1(tx.Height, proposalID)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	var sb strings.Builder

	for i, proposalMsg := range proposal.Messages {
		if i == 0 {
			sb.WriteString("[")
		}

		bz, err := m.db.EncodingConfig.Codec.MarshalJSON(proposalMsg)
		if err != nil {
			return fmt.Errorf("error while marshaling proposal content: %s", err)
		}

		sb.Write(bz)

		if (i + 1) == len(proposal.Messages) {
			sb.WriteString("]")
		} else {
			sb.WriteString(",")
		}
	}

	proposalObj := types.NewProposal(
		proposal.Id,
		sb.String(),
		proposal.Status.String(),
		proposal.SubmitBlock,
		proposal.DepositEndBlock,
		proposal.VotingStartBlock,
		proposal.VotingEndBlock,
		msg.Proposer,
		proposal.Metadata,
	)
	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.Id, msg.Proposer, msg.InitialDeposit, tx.Height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func (m *Module) handleMsgSubmitProposal(tx *juno.Tx, index int, msg *govtypesv1beta1.MsgSubmitProposal) error {
	proposalID, err := m.getProposalID(tx, index)
	if err != nil {
		return fmt.Errorf("error while getting proposal id: %s", err)
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

	// Encode the content properly
	protoContent, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid proposal content types: %T", proposal.Content)
	}

	anyContent, err := codectypes.NewAnyWithValue(protoContent)
	if err != nil {
		return fmt.Errorf("error while wrapping proposal proto content: %s", err)
	}

	contentBz, err := m.db.EncodingConfig.Codec.MarshalJSON(anyContent)
	if err != nil {
		return fmt.Errorf("error while marshaling proposal content: %s", err)
	}

	metadata := map[string]string{
		"title":       proposal.GetContent().GetTitle(),
		"description": proposal.GetContent().GetDescription(),
		"type":        proposal.ProposalType(),
	}

	metadataBz, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error while marshaling proposal metadata: %s", err)
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.ProposalId,
		string(contentBz),
		proposal.Status.String(),
		proposal.SubmitBlock,
		proposal.DepositEndBlock,
		proposal.VotingStartBlock,
		proposal.VotingEndBlock,
		msg.Proposer,
		string(metadataBz),
	)
	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

func (m *Module) getProposalID(tx *juno.Tx, index int) (uint64, error) {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return 0, fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return 0, fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error while parsing proposal id: %s", err)
	}

	return proposalID, nil
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func (m *Module) handleMsgDeposit(height int64, proposalId uint64, depositor string) error {
	deposit, err := m.source.ProposalDeposit(height, proposalId, depositor)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	return m.db.SaveDeposits([]types.Deposit{
		types.NewDeposit(proposalId, depositor, deposit.Amount, height),
	})
}

// handleMsgVote allows to properly handle a handleMsgVote
func (m *Module) handleMsgVote(height int64, proposalId uint64, voter string, option int32) error {
	vote := types.NewVote(proposalId, voter, option, height)
	return m.db.SaveVote(vote)
}

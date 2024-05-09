package gov

import (
	"encoding/json"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	juno "github.com/forbole/juno/v4/types"
	"github.com/gogo/protobuf/proto"
	"github.com/rarimo/bdjuno/types"
	"github.com/rarimo/rarimo-core/app"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	"strconv"

	"github.com/rs/zerolog/log"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	if err := m.handleBlockEvents(b.Block.Height, res.BeginBlockEvents); err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while processing end block events")
	}

	if err := m.handleBlockEvents(b.Block.Height, res.EndBlockEvents); err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while processing end block events")
	}

	err := m.updateProposals(b.Block.Height, vals)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating proposals")
	}
	return nil
}

// updateProposals updates the proposals
func (m *Module) updateProposals(height int64, blockVals *tmctypes.ResultValidators) error {
	ids, err := m.db.GetOpenProposalsIds(uint64(height))
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open ids")
	}

	for _, id := range ids {
		err = m.UpdateProposal(height, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal: %s", err)
		}

		err = m.UpdateProposalSnapshots(height, blockVals, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal snapshots: %s", err)
		}
	}
	return nil
}

func (m *Module) handleBlockEvents(height int64, events []abci.Event) error {
	for _, event := range events {
		switch event.Type {
		case govtypes.EventTypeSubmitProposal:
			if err := m.handleBlockEventSubmitProposal(height, event); err != nil {
				return err
			}
			// Add other events handling if required
		}
	}

	return nil
}

func (m *Module) handleBlockEventSubmitProposal(height int64, event abci.Event) error {
	proposalIdAttribute, err := juno.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		// error means to such attribute - normal logic
		return nil
	}

	proposalId, err := strconv.ParseUint(string(proposalIdAttribute.Value), 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing proposal id atribute: %s", err)
	}

	// !! Logic partially copied from handleMsgSubmitProposal method

	proposal, err := m.source.Proposal(height, proposalId)
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

	// !! If proposal was created in Block events then the proposer will be defined as module account
	govModuleAddress, err := bech32.ConvertAndEncode(
		app.AccountAddressPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).Bytes(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to convert module address %s", err))
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
		govModuleAddress,
		string(metadataBz),
	)
	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, govModuleAddress, proposal.TotalDeposit, height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

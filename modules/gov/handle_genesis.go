package gov

import (
	"encoding/json"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"
	tmtypes "github.com/tendermint/tendermint/types"

	"gitlab.com/rarimo/bdjuno/types"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genStatev1beta1 govtypesv1beta1.GenesisState
	err := m.cdc.UnmarshalJSON(appState[gov.ModuleName], &genStatev1beta1)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// Save the proposals
	err = m.saveGenesisProposals(genStatev1beta1.Proposals, doc)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	// Save the params
	err = m.db.SaveGenesisGovParams(types.NewGenesisGovParams(
		types.NewGenesisVotingParams(&genStatev1beta1.VotingParams),
		types.NewGenesisDepositParam(&genStatev1beta1.DepositParams),
		types.NewGenesisTallyParams(&genStatev1beta1.TallyParams),
		doc.InitialHeight,
	))
	if err != nil {
		return fmt.Errorf("error while storing genesis governance params: %s", err)
	}

	return nil
}

// / saveGenesisProposals save proposals from genesis file
func (m *Module) saveGenesisProposals(slice govtypesv1beta1.Proposals, genDoc *tmtypes.GenesisDoc) error {
	proposals := make([]types.Proposal, len(slice))
	tallyResults := make([]types.TallyResult, len(slice))
	deposits := make([]types.Deposit, len(slice))

	for index, proposal := range slice {
		// Since it's not possible to get the proposer, set it to nil

		// Unpack the content
		var content govtypesv1beta1.Content
		err := m.cdc.UnpackAny(proposal.Content, &content)
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

		proposals[index] = types.NewProposal(
			proposal.ProposalId,
			string(contentBz),
			proposal.Status.String(),
			proposal.SubmitBlock,
			proposal.DepositEndBlock,
			proposal.VotingStartBlock,
			proposal.VotingEndBlock,
			"",
			string(metadataBz),
		)

		tallyResults[index] = types.NewTallyResult(
			proposal.ProposalId,
			proposal.FinalTallyResult.Yes.String(),
			proposal.FinalTallyResult.Abstain.String(),
			proposal.FinalTallyResult.No.String(),
			proposal.FinalTallyResult.NoWithVeto.String(),
			genDoc.InitialHeight,
		)

		deposits[index] = types.NewDeposit(
			proposal.ProposalId,
			"",
			proposal.TotalDeposit,
			genDoc.InitialHeight,
		)
	}

	// Save the proposals
	err := m.db.SaveProposals(proposals)
	if err != nil {
		return err
	}

	// Save the deposits
	err = m.db.SaveDeposits(deposits)
	if err != nil {
		return err
	}

	// Save the tally results
	return m.db.SaveTallyResults(tallyResults)
}

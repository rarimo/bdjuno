package multisig

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v4/types"
	multisigtypes "gitlab.com/rarimo/rarimo-core/x/multisig/types"
	"strconv"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *multisigtypes.MsgCreateGroup:
		return m.handleMsgCreateGroup(tx, index)
	case *multisigtypes.MsgChangeGroup:
		return m.handleMsgChangeGroup(tx.Height, cosmosMsg)
	case *multisigtypes.MsgSubmitProposal:
		return m.handleMsgSubmitProposal(tx, index)
	case *multisigtypes.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)
	}

	return nil
}

func (m *Module) handleMsgCreateGroup(tx *juno.Tx, index int) error {
	// Get the group id
	event, err := tx.FindEventByType(index, multisigtypes.EventTypeCreateGroup)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeCreateGroup: %s", err)
	}

	account, err := tx.FindAttributeByKey(event, multisigtypes.AttributeKeyGroup)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyGroup: %s", err)
	}

	return m.saveGroup(tx.Height, account)
}

func (m *Module) handleMsgChangeGroup(height int64, msg *multisigtypes.MsgChangeGroup) error {
	return m.saveGroup(height, msg.Group)
}

func (m *Module) handleMsgSubmitProposal(tx *juno.Tx, index int) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, multisigtypes.EventTypeSubmitProposal)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	raw, err := tx.FindAttributeByKey(event, multisigtypes.AttributeKeyProposal)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyProposal: %s", err)
	}

	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing AttributeKeyProposal: %s", err)
	}

	return m.saveProposal(tx.Height, id)
}

func (m *Module) handleMsgVote(tx *juno.Tx, msg *multisigtypes.MsgVote) error {
	return m.saveVote(tx.Height, msg.ProposalId, msg.Creator)
}

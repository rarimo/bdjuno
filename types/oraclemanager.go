package types

import (
	oraclemanagertypes "github.com/rarimo/rarimo-core/x/oraclemanager/types"
)

// OracleManagerParams contains the data of the x/oraclemanager module params instance
type OracleManagerParams struct {
	MinOracleStake      string `json:"min_oracle_stake,omitempty" yaml:"min_oracle_stake,omitempty"`
	CheckOperationDelta uint64 `json:"check_operation_delta,omitempty" yaml:"check_operation_delta,omitempty"`
	MaxViolationsCount  uint64 `json:"max_violations_count,omitempty" yaml:"max_violations_count,omitempty"`
	MaxMissedCount      uint64 `json:"max_missed_count,omitempty" yaml:"max_missed_count,omitempty"`
	SlashedFreezeBlocks uint64 `json:"slashed_freeze_blocks,omitempty" yaml:"slashed_freeze_blocks,omitempty"`
	MinOraclesCount     uint64 `json:"min_oracles_count,omitempty" yaml:"min_oracles_count,omitempty"`
	StakeDenom          string `json:"stake_denom,omitempty" yaml:"stake_denom,omitempty"`
	VoteQuorum          string `json:"vote_quorum,omitempty" yaml:"vote_quorum,omitempty"`
	VoteThreshold       string `json:"vote_threshold,omitempty" yaml:"vote_threshold,omitempty"`
	Height              int64  `json:"height,omitempty" yaml:"height,omitempty"`
}

// OracleManagerParamsFromCore allows to build a new OracleManagerParams instance from an oraclemanagertypes.Params instance
func OracleManagerParamsFromCore(p oraclemanagertypes.Params, height int64) *OracleManagerParams {
	return &OracleManagerParams{
		MinOracleStake:      p.MinOracleStake,
		CheckOperationDelta: p.CheckOperationDelta,
		MaxViolationsCount:  p.MaxViolationsCount,
		MaxMissedCount:      p.MaxMissedCount,
		SlashedFreezeBlocks: p.SlashedFreezeBlocks,
		MinOraclesCount:     p.MinOraclesCount,
		StakeDenom:          p.StakeDenom,
		VoteQuorum:          p.VoteQuorum,
		VoteThreshold:       p.VoteThreshold,
		Height:              height,
	}
}

// Oracle represents a single oracle instance
type Oracle struct {
	Index                 string                          `json:"index,omitempty" yaml:"index,omitempty"`
	Chain                 string                          `json:"chain,omitempty" yaml:"chain,omitempty"`
	Account               string                          `json:"account,omitempty" yaml:"account,omitempty"`
	Status                oraclemanagertypes.OracleStatus `json:"status,omitempty" yaml:"status,omitempty"`
	Stake                 string                          `json:"stake,omitempty" yaml:"stake,omitempty"`
	MissedCount           uint64                          `json:"missed_count,omitempty" yaml:"missed_count,omitempty"`
	ViolationsCount       uint64                          `json:"violations_count,omitempty" yaml:"violations_count,omitempty"`
	FreezeEndBlock        uint64                          `json:"freeze_end_block,omitempty" yaml:"freeze_end_block,omitempty"`
	VotesCount            uint64                          `json:"votes_count,omitempty" yaml:"votes_count,omitempty"`
	CreateOperationsCount uint64                          `json:"create_operations_count,omitempty" yaml:"create_operations_count,omitempty"`
}

// OracleFromCore allows to build a new Oracle instance from an oraclemanagertypes.Oracle instance
func OracleFromCore(oracle oraclemanagertypes.Oracle) Oracle {
	return Oracle{
		Index:                 string(oraclemanagertypes.OracleKey(oracle.Index)),
		Chain:                 oracle.Index.Chain,
		Account:               oracle.Index.Account,
		Status:                oracle.Status,
		Stake:                 oracle.Stake,
		MissedCount:           oracle.MissedCount,
		ViolationsCount:       oracle.ViolationsCount,
		FreezeEndBlock:        oracle.FreezeEndBlock,
		VotesCount:            oracle.VotesCount,
		CreateOperationsCount: oracle.CreateOperationsCount,
	}
}

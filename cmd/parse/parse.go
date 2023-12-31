package parse

import (
	parse "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"

	parseblocks "github.com/forbole/juno/v4/cmd/parse/blocks"

	parsegenesis "github.com/forbole/juno/v4/cmd/parse/genesis"

	parsetransaction "github.com/forbole/juno/v4/cmd/parse/transactions"
	parseauth "github.com/rarimo/bdjuno/cmd/parse/auth"
	parsefeegrant "github.com/rarimo/bdjuno/cmd/parse/feegrant"
	parsegov "github.com/rarimo/bdjuno/cmd/parse/gov"
	parsestaking "github.com/rarimo/bdjuno/cmd/parse/staking"
)

// NewParseCmd returns the Cobra command allowing to parse some chain data without having to re-sync the whole database
func NewParseCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "parse",
		Short:             "Parse some data without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: runPersistentPreRuns(parse.ReadConfigPreRunE(parseCfg)),
	}

	cmd.AddCommand(
		parseauth.NewAuthCmd(parseCfg),
		parseblocks.NewBlocksCmd(parseCfg),
		parsefeegrant.NewFeegrantCmd(parseCfg),
		parsegenesis.NewGenesisCmd(parseCfg),
		parsegov.NewGovCmd(parseCfg),
		parsestaking.NewStakingCmd(parseCfg),
		parsetransaction.NewTransactionsCmd(parseCfg),
	)

	return cmd
}

func runPersistentPreRuns(preRun func(_ *cobra.Command, _ []string) error) func(_ *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if root := cmd.Root(); root != nil {
			if root.PersistentPreRunE != nil {
				err := root.PersistentPreRunE(root, args)
				if err != nil {
					return err
				}
			}
		}

		return preRun(cmd, args)
	}
}

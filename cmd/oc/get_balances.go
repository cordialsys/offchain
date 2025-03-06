package main

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/okx"
	"github.com/spf13/cobra"
)

func NewGetBalancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "balances",
		Short:        "Get the balances of the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := okx.NewClient(loadFromEnv())
			if err != nil {
				return err
			}

			assets, err := cli.GetBalances(oc.GetBalanceArgs{})
			if err != nil {
				return err
			}
			printJson(assets)
			return nil
		},
	}
	return cmd
}

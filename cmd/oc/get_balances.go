package main

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewGetBalancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "balances",
		Short:        "Get the balances of the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig := unwrapExchangeConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig)
			if err != nil {
				return err
			}

			assets, err := cli.ListBalances(oc.GetBalanceArgs{})
			if err != nil {
				return err
			}
			printJson(assets)
			return nil
		},
	}
	return cmd
}

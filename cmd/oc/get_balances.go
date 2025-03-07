package main

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewGetBalancesCmd() *cobra.Command {
	var account string
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
			accountName := oc.CoreFunding
			if account != "" {
				accountName = oc.NewAccountName(account)
			}
			balanceArgs := oc.NewGetBalanceArgs(accountName)

			assets, err := cli.ListBalances(balanceArgs)
			if err != nil {
				return err
			}
			printJson(assets)
			return nil
		},
	}
	cmd.Flags().StringVar(
		&account,
		"account",
		"",
		"the account to get balances for",
	)
	return cmd
}

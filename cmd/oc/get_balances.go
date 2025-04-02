package main

import (
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewListBalancesCmd() *cobra.Command {
	var account string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "balances",
		Short:        "List your balances on the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig := unwrapExchangeConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig)
			if err != nil {
				return err
			}
			accountName := client.CoreFunding
			if account != "" {
				accountName = client.NewAccountName(account)
			}
			balanceArgs := client.NewGetBalanceArgs(accountName)

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

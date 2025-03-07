package main

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewGetDepositAddressCmd() *cobra.Command {
	var symbol string
	var network string
	var account string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "deposit",
		Short:        "Get a deposit address for a symbol and network",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig := unwrapExchangeConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig)
			if err != nil {
				return err
			}
			if symbol == "" {
				return fmt.Errorf("--symbol is required")
			}
			if network == "" {
				return fmt.Errorf("--network is required")
			}

			options := []oc.GetDepositAddressOption{}
			if account != "" {
				options = append(options, oc.WithAccount(oc.AccountName(account)))
			}

			resp, err := cli.GetDepositAddress(oc.NewGetDepositAddressArgs(
				oc.SymbolId(symbol),
				oc.NetworkId(network),
				options...,
			))

			if err != nil {
				return err
			}
			fmt.Println(resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to withdraw")
	cmd.Flags().StringVar(&network, "network", "", "The network to transact on")
	cmd.Flags().StringVar(&account, "account", "", "The account to get the deposit address for (optional)")
	return cmd
}

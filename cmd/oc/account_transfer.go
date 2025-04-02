package main

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewAccountTransferCmd() *cobra.Command {
	var from string
	var to string
	var symbol string
	var amountS string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "transfer",
		Short:        "Transfer funds between accounts on exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig := unwrapExchangeConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig)
			if err != nil {
				return err
			}
			amount, err := oc.NewAmountFromString(amountS)
			if err != nil {
				return err
			}
			if amount.IsZero() {
				return fmt.Errorf("--amount must be greater than 0")
			}
			if symbol == "" {
				return fmt.Errorf("--symbol is required")
			}

			resp, err := cli.CreateAccountTransfer(client.NewAccountTransferArgs(
				client.AccountName(from),
				client.AccountName(to),
				oc.SymbolId(symbol),
				amount,
			))
			if err != nil {
				return err
			}
			printJson(resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", string(client.CoreFunding), "The account to transfer from")
	cmd.Flags().StringVar(&to, "to", string(client.CoreTrading), "The account to transfer to")
	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to transfer")
	cmd.Flags().StringVar(&amountS, "amount", "", "The amount to transfer")
	return cmd
}

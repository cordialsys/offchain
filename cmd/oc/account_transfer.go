package main

import (
	oc "github.com/cordialsys/offchain"
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

			err = cli.CreateAccountTransfer(oc.NewAccountTransferArgs(
				oc.AccountName(from),
				oc.AccountName(to),
				oc.SymbolId(symbol),
				amount,
			))
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", string(oc.CoreFunding), "The account to transfer from")
	cmd.Flags().StringVar(&to, "to", string(oc.CoreTrading), "The account to transfer to")
	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to transfer")
	cmd.Flags().StringVar(&amountS, "amount", "", "The amount to transfer")
	return cmd
}

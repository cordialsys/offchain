package main

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/exchanges/okx"
	"github.com/cordialsys/offchain/exchanges/okx/api"
	"github.com/spf13/cobra"
)

const InternalTransfer = "3"
const OnChainTransfer = "4"

func NewWithdrawCmd() *cobra.Command {
	var from string
	var to string
	var symbol string
	var network string
	var amountS string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "withdraw",
		Short:        "Withdraw funds from the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := okx.NewClient(loadFromEnv())
			if err != nil {
				return err
			}
			apiokx := cli.Api()
			amount, err := oc.NewAmountFromString(amountS)
			if err != nil {
				return err
			}

			response, err := apiokx.Withdrawal(&api.WithdrawalRequest{
				Amount:      amount.String(),
				Destination: OnChainTransfer,
				Currency:    symbol,
				Chain:       network,
				ToAddress:   to,
			})
			printJson(response)

			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", string(oc.CoreFunding), "The account to transfer from")
	cmd.Flags().StringVar(&to, "to", "", "Your address to withdraw to")
	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to withdraw")
	cmd.Flags().StringVar(&symbol, "network", "", "The network to transact on")
	cmd.Flags().StringVar(&amountS, "amount", "", "The amount to withdraw")
	return cmd
}

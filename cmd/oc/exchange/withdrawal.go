package exchange

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/spf13/cobra"
)

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
			cli := unwrapClient(cmd.Context())
			amount, err := oc.NewAmountFromString(amountS)
			if err != nil {
				return err
			}

			resp, err := cli.CreateWithdrawal(client.NewWithdrawalArgs(
				oc.Address(to),
				oc.SymbolId(symbol),
				oc.NetworkId(network),
				amount,
			))

			if err != nil {
				return err
			}
			printJson(resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", string(""), "The account to transfer from")
	cmd.Flags().StringVar(&to, "to", "", "Your address to withdraw to")
	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to withdraw")
	cmd.Flags().StringVar(&network, "network", "", "The network to transact on")
	cmd.Flags().StringVar(&amountS, "amount", "", "The amount to withdraw")
	return cmd
}

package exchange

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/spf13/cobra"
)

func NewGetDepositAddressCmd() *cobra.Command {
	var symbol string
	var network string
	var subaccount string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "deposit",
		Short:        "Get a deposit address for a symbol and network",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := unwrapClient(cmd.Context())
			if symbol == "" {
				return fmt.Errorf("--symbol is required")
			}
			if network == "" {
				return fmt.Errorf("--network is required")
			}

			options := []client.GetDepositAddressOption{}
			if subaccount != "" {
				options = append(options, client.WithSubaccount(oc.AccountId(subaccount)))
			}

			resp, err := cli.GetDepositAddress(client.NewGetDepositAddressArgs(
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
	cmd.Flags().StringVar(&subaccount, "for", "", "The subaccount to get the deposit address for, when using the main account to query (optional)")
	return cmd
}

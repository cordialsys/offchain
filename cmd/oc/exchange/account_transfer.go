package exchange

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
	var fromType string
	var toType string

	var symbol string
	var amountS string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "transfer",
		Short:        "Transfer funds between accounts on exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig, secrets := unwrapAccountConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig.ExchangeId, &exchangeConfig.ExchangeClientConfig, secrets)
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

			transferArgs := client.NewAccountTransferArgs(
				// client.AccountName(from),
				// client.AccountName(to),
				oc.SymbolId(symbol),
				amount,
			)
			transferArgs.SetFrom(client.AccountName(from), client.AccountType(fromType))
			transferArgs.SetTo(client.AccountName(to), client.AccountType(toType))

			resp, err := cli.CreateAccountTransfer(transferArgs)
			if err != nil {
				return err
			}

			printJson(resp)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "The account to transfer from (defaults to the main account)")
	cmd.Flags().StringVar(&to, "to", "", "The account to transfer to (defaults to the main account)")
	cmd.Flags().StringVar(&fromType, "from-type", "", "The type of account to transfer from (defaults to core-funding)")
	cmd.Flags().StringVar(&toType, "to-type", "", "The type of account to transfer to (defaults to core-trading)")

	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to transfer")
	cmd.Flags().StringVar(&amountS, "amount", "", "The amount to transfer")
	return cmd
}

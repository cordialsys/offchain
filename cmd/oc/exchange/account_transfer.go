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
			cli, err := loader.NewClient(exchangeConfig, secrets)
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
			var message string
			transferArgs := client.NewAccountTransferArgs(
				oc.SymbolId(symbol),
				amount,
			)

			if fromType != "" {
				resolvedFromType, ok, message := exchangeConfig.ResolveAccountType(fromType)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				transferArgs.SetFromType(resolvedFromType.Type)
			} else {
				// use the first account type by default
				first, ok := exchangeConfig.FirstAccountType()
				if ok {
					transferArgs.SetFromType(first.Type)
				}
			}
			if toType != "" {
				resolvedToType, ok, message := exchangeConfig.ResolveAccountType(toType)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				transferArgs.SetToType(resolvedToType.Type)
			} else {
				// use the first account type by default
				first, ok := exchangeConfig.FirstAccountType()
				if ok {
					transferArgs.SetToType(first.Type)
				}
			}

			if from != "" {
				fromSubaccount, ok := exchangeConfig.ResolveSubAccount(from)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				transferArgs.SetFrom(fromSubaccount.Id)
			}
			if to != "" {
				toSubaccount, ok := exchangeConfig.ResolveSubAccount(to)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				transferArgs.SetTo(toSubaccount.Id)
			}

			if (to + from + fromType + toType) == "" {
				return fmt.Errorf("must specify at least one of --to, --from, --from-type, --to-type")
			}

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
	cmd.Flags().StringVar(&fromType, "from-type", "", "The type of account to transfer from (defaults to first account type)")
	cmd.Flags().StringVar(&toType, "to-type", "", "The type of account to transfer to (defaults to first account type)")

	cmd.Flags().StringVar(&symbol, "symbol", "", "The symbol to transfer")
	cmd.Flags().StringVar(&amountS, "amount", "", "The amount to transfer")
	return cmd
}

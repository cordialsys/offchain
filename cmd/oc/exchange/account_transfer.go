package exchange

import (
	"fmt"
	"slices"
	"strings"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func resolveAccountType(cli loader.Client, typeOrAlias string) (*oc.AccountTypeConfig, bool, error) {
	accountTypes, err := cli.ListAccountTypes()
	options := []string{}
	if err != nil {
		return nil, false, fmt.Errorf("failed to list account types: %s", err)
	}
	for _, accountType := range accountTypes {
		if string(accountType.Type) == typeOrAlias || slices.Contains(accountType.Aliases, typeOrAlias) {
			return accountType, true, nil
		}
		if len(accountType.Aliases) > 0 {
			options = append(options, fmt.Sprintf("%s (%s)", accountType.Type, accountType.Aliases[0]))
		} else {
			options = append(options, string(accountType.Type))
		}
	}
	return nil, false, fmt.Errorf("account type or alias %s not found; options: %s", typeOrAlias, strings.Join(options, ", "))
}

func resolveFirstAccountType(cli loader.Client) (*oc.AccountTypeConfig, bool, error) {
	accountTypes, err := cli.ListAccountTypes()
	if err != nil {
		return nil, false, fmt.Errorf("failed to list account types: %s", err)
	}
	if len(accountTypes) == 0 {
		return nil, false, nil
	}
	return accountTypes[0], true, nil
}

func resolveSubAccount(cli loader.Client, idOrAlias string) (accountCfg *oc.SubAccountHeader, err error) {
	subaccounts, err := cli.ListSubaccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to list subaccounts: %s", err)
	}
	for _, sa := range subaccounts {
		if string(sa.Id) == idOrAlias || sa.Alias == idOrAlias {
			return sa, nil
		}
	}
	return nil, fmt.Errorf("subaccount %s not found", idOrAlias)
}

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
			cli := unwrapClient(cmd.Context())
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
				oc.SymbolId(symbol),
				amount,
			)

			if fromType != "" {

				accountTypes, err := cli.ListAccountTypes()
				if err != nil {
					return fmt.Errorf("failed to list account types: %s", err)
				}
				for _, accountType := range accountTypes {
					if string(accountType.Type) == fromType || slices.Contains(accountType.Aliases, fromType) {
					}
				}

				resolvedFromType, ok, message := resolveAccountType(cli, fromType)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				transferArgs.SetFromType(resolvedFromType.Type)
			} else {
				// use the first account type by default
				first, ok, err := resolveFirstAccountType(cli)
				if err != nil {
					return err
				}
				if ok {
					transferArgs.SetFromType(first.Type)
				}
			}
			if toType != "" {
				resolvedToType, ok, message := resolveAccountType(cli, toType)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				transferArgs.SetToType(resolvedToType.Type)
			} else {
				// use the first account type by default
				first, ok, err := resolveFirstAccountType(cli)
				if err != nil {
					return err
				}
				if ok {
					transferArgs.SetToType(first.Type)
				}
			}

			if from != "" {
				fromSubaccount, err := resolveSubAccount(cli, from)
				if err != nil {
					return err
				}
				transferArgs.SetFrom(fromSubaccount.Id)
			}
			if to != "" {
				toSubaccount, err := resolveSubAccount(cli, to)
				if err != nil {
					return err
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

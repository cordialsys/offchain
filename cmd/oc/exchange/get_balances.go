package exchange

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewListBalancesCmd() *cobra.Command {
	var accountType string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "balances",
		Short:        "List your balances on the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig, secrets := unwrapAccountConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig, secrets)
			if err != nil {
				return err
			}
			balanceArgs := client.NewGetBalanceArgs("")
			if accountType != "" {
				at, ok, message := exchangeConfig.ResolveAccountType(accountType)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				balanceArgs.SetAccountType(at.Type)
			} else {
				// try taking the first account type
				defaults, ok := oc.GetDefaultConfig(exchangeConfig.ExchangeId)
				if ok {
					if len(defaults.AccountTypes) > 0 {
						logrus.WithFields(logrus.Fields{
							"type": defaults.AccountTypes[0].Type,
						}).Infof("using default account type")
						balanceArgs.SetAccountType(defaults.AccountTypes[0].Type)
					}
				}
			}

			assets, err := cli.ListBalances(balanceArgs)
			if err != nil {
				return err
			}
			printJson(assets)
			return nil
		},
	}
	cmd.Flags().StringVar(
		&accountType,
		"type",
		"",
		"the account type to get balances for",
	)
	return cmd
}

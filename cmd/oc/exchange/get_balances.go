package exchange

import (
	"fmt"

	"github.com/cordialsys/offchain/client"
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
			cli := unwrapClient(cmd.Context())
			balanceArgs := client.NewGetBalanceArgs("")
			if accountType != "" {
				at, ok, message := resolveAccountType(cli, accountType)
				if !ok {
					return fmt.Errorf("%s", message)
				}
				balanceArgs.SetAccountType(at.Type)
			} else {
				// try taking the first account type
				at, ok, err := resolveFirstAccountType(cli)
				if err != nil {
					return err
				}
				if ok {
					logrus.WithFields(logrus.Fields{
						"type": at.Type,
					}).Infof("using default account type")
					balanceArgs.SetAccountType(at.Type)
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

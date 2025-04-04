package exchange

import (
	"fmt"

	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewListWithdrawalHistoryCmd() *cobra.Command {
	var limit int
	var pageToken string
	cmd := &cobra.Command{
		Use:   "history",
		Short: "List withdrawal history",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig, secrets := unwrapAccountConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig.ExchangeId, &exchangeConfig.ExchangeClientConfig, secrets)
			if err != nil {
				return err
			}

			if cmd.Flags().Changed("limit") {
				return fmt.Errorf("--limit is not yet supported")
			}
			if cmd.Flags().Changed("page-token") {
				return fmt.Errorf("--page-token is not yet supported")
			}

			resp, err := cli.ListWithdrawalHistory(client.NewWithdrawalHistoryArgs())
			if err != nil {
				return err
			}
			printJson(resp)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 100, "The number of items to return")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Page token to use")
	return cmd
}

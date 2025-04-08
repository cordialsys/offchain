package exchange

import (
	"github.com/spf13/cobra"
)

// This just reads the configuration currently, not yet exposing an API call for this.
func NewListSubaccountsCmd() *cobra.Command {
	return &cobra.Command{
		SilenceUsage: true,
		Use:          "subaccounts",
		Short:        "List all configured subaccounts on an exchange",

		RunE: func(cmd *cobra.Command, args []string) error {
			cli := unwrapClient(cmd.Context())
			subaccounts, err := cli.ListSubaccounts()
			if err != nil {
				return err
			}
			printJson(subaccounts)
			return nil
		},
	}
}

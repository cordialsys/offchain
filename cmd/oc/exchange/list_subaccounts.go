package exchange

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

// This just reads the configuration currently, not yet exposing an API call for this.
func NewListSubaccountsCmd() *cobra.Command {
	var exchangeConfig *oc.ExchangeConfig
	return &cobra.Command{
		SilenceUsage: true,
		Use:          "list-subaccounts",
		Aliases:      []string{"subaccounts"},
		Short:        "List all configured subaccounts on an exchange",
		// override the default pre-run to only load the config
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			exchange := cmd.Flag("exchange").Value.String()
			configPath := cmd.Flag("config").Value.String()
			config, err := loader.LoadConfig(configPath)
			if err != nil {
				return err
			}
			var ok bool

			exchangeConfig, ok = config.GetExchange(oc.ExchangeId(exchange))
			if !ok {
				return fmt.Errorf("exchange not found. options are: %v", oc.ValidExchangeIds)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			subaccounts := []string{}
			for _, subaccount := range exchangeConfig.SubAccounts {
				subaccounts = append(subaccounts, subaccount.Id)
			}
			printJson(subaccounts)
			return nil
		},
	}
}

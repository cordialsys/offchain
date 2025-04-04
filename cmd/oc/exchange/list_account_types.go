package exchange

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/spf13/cobra"
)

var NopSecrets = &oc.MultiSecret{
	ApiKeyRef:     secret.Secret("raw:AAA"),
	SecretKeyRef:  secret.Secret("raw:AAA"),
	PassphraseRef: secret.Secret("raw:AAA"),
}

// This just reads the configuration currently, not yet exposing an API call for this.
func NewListAccountTypesCmd() *cobra.Command {
	var exchangeConfig *oc.ExchangeConfig
	return &cobra.Command{
		SilenceUsage: true,
		Use:          "account-types",
		Short:        "List valid account types for an exchange",
		// override the default pre-run to only load the config
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			exchange := cmd.Flag("exchange").Value.String()
			configPath := cmd.Flag("config").Value.String()
			config, err := loader.LoadValidatedConfig(configPath)
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
			cli, err := loader.NewClient(exchangeConfig, NopSecrets)
			if err != nil {
				return err
			}
			accountTypes, err := cli.ListAccountTypes()
			if err != nil {
				return err
			}
			printJson(accountTypes)
			return nil
		},
	}
}

package exchange

import (
	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
)

func NewGetAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "assets",
		Short:        "list the asset symbols and networks of the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			exchangeConfig, secrets := unwrapAccountConfig(cmd.Context())
			cli, err := loader.NewClient(exchangeConfig.ExchangeId, &exchangeConfig.ExchangeClientConfig, secrets)
			if err != nil {
				return err
			}
			assets, err := cli.ListAssets()
			if err != nil {
				return err
			}
			printJson(assets)
			return nil
		},
	}
	return cmd
}

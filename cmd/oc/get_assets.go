package main

import (
	"github.com/cordialsys/offchain/exchanges/okx"
	"github.com/spf13/cobra"
)

func NewGetAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "assets",
		Short:        "Get the asset symbols and networks of the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {

			cli, err := okx.NewClient(loadFromEnv())
			if err != nil {
				return err
			}

			assets, err := cli.GetAssets()
			if err != nil {
				return err
			}
			printJson(assets)
			return nil
		},
	}
	return cmd
}

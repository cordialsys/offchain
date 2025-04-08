package exchange

import (
	"github.com/spf13/cobra"
)

func NewGetAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "assets",
		Short:        "list the asset symbols and networks of the exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli := unwrapClient(cmd.Context())
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

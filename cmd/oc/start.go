package main

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	var listen string
	var configPath string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Serve the offchain server",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loader.LoadValidatedConfig(configPath)
			if err != nil {
				return err
			}
			server := server.New(listen, config)
			return server.Start()
		},
	}
	cmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		fmt.Sprintf("path to the config file (may set %s)", oc.ENV_OFFCHAIN_CONFIG),
	)
	cmd.Flags().StringVarP(&listen, "listen", "l", ":8080", "address to listen on")

	return cmd
}

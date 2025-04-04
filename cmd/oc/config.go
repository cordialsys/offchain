package main

import (
	"fmt"

	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewConfigCmd() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "config",
		Short:        "Print out the current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loader.LoadConfig(configPath)
			if err != nil {
				return err
			}
			out, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().StringVarP(
		&configPath,
		"config",
		"c",
		"",
		fmt.Sprintf("path to the config file (may set %s)", loader.ENV_OFFCHAIN_CONFIG),
	)
	return cmd
}

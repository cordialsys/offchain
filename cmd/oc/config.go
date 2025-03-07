package main

import (
	"fmt"

	"github.com/cordialsys/offchain/loader"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "config",
		Short:        "Print out the current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cmd.Context().Value(configContextKey).(*loader.Config)
			out, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
	return cmd
}

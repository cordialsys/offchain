package main

import (
	"fmt"

	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/spf13/cobra"
)

func NewSecretCmd() *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage: true,
		Use:          "secret <reference>",
		Short:        "Print out a secret given a reference",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			sec, err := secret.Secret(args[0]).Load()
			if err != nil {
				return err
			}
			fmt.Println(sec)
			return nil
		},
	}
	return cmd
}
